// Package main is the entry point for the TeslaPay API Gateway.
// It sets up HTTP routing, middleware, and coordinates requests to
// downstream microservices.
//
// In production, Kong is the external API Gateway. This Go gateway
// serves as the internal routing layer behind Kong, handling service
// composition and providing a unified HTTP interface to the Flutter client.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/account"
	"github.com/teslapay/backend/internal/auth"
	"github.com/teslapay/backend/internal/card"
	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/internal/crypto"
	"github.com/teslapay/backend/internal/kyc"
	"github.com/teslapay/backend/internal/notification"
	"github.com/teslapay/backend/internal/payment"
	tpcrypto "github.com/teslapay/backend/pkg/crypto"
	"github.com/teslapay/backend/pkg/database"
	"github.com/teslapay/backend/pkg/events"
	"github.com/teslapay/backend/pkg/middleware"
)

func main() {
	// Load configuration from environment variables.
	cfg, err := common.LoadBaseConfig("api-gateway")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize structured logger.
	logger, err := common.NewLogger(cfg.ServiceName, cfg.Environment, cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("starting TeslaPay API Gateway",
		zap.String("environment", cfg.Environment),
		zap.Int("port", cfg.Port),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize JWT manager. In development, keys may not be present.
	jwtManager, err := initJWTManager(cfg, logger)
	if err != nil {
		logger.Warn("JWT manager initialization failed, auth endpoints will be limited",
			zap.Error(err),
		)
	}

	// Initialize Redis client for rate limiting and sessions.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn("Redis connection failed, rate limiting will be disabled", zap.Error(err))
	}

	// Initialize database connection (for the gateway's own needs).
	var db *database.DB
	if cfg.DatabaseURL != "" {
		dbCfg := database.DefaultConfig(cfg.DatabaseURL)
		db, err = database.New(ctx, dbCfg, logger)
		if err != nil {
			logger.Error("database connection failed", zap.Error(err))
			// Gateway can run without direct DB access in some modes.
		}
	}

	// Set up Gin router.
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Global middleware.
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLogger(logger))

	// Rate limiter (gracefully degrades if Redis is down).
	rateLimiter := middleware.NewRateLimiter(redisClient, logger)

	// Health check endpoint (no auth required).
	router.GET("/health", healthCheck(db, redisClient))
	router.GET("/ready", readinessCheck(db))

	// Public API v1.
	v1 := router.Group("/api/v1")
	v1.Use(rateLimiter.Middleware())

	// Auth routes (mostly unauthenticated).
	if db != nil {
		authRepo := auth.NewRepository(db, logger)
		authService := auth.NewService(authRepo, jwtManager, logger, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
		authHandler := auth.NewHandler(authService, logger)
		authHandler.RegisterRoutes(v1)
	}

	// Protected routes (require JWT authentication).
	protected := v1.Group("")
	if jwtManager != nil {
		protected.Use(middleware.AuthMiddleware(jwtManager, logger))
	}

	if db != nil {
		accountRepo := account.NewRepository(db, logger)
		accountService := account.NewService(accountRepo, logger)
		accountHandler := account.NewHandler(accountService, logger)
		accountHandler.RegisterRoutes(protected)
	}

	if db != nil {
		// Initialize Kafka producer for gateway (used by payment and card services).
		var sharedProducer *events.Producer
		if len(cfg.KafkaBrokers) > 0 {
			producerCfg := events.DefaultProducerConfig(cfg.KafkaBrokers)
			sharedProducer = events.NewProducer(producerCfg, logger)
		}

		paymentRepo := payment.NewRepository(db, logger)
		paymentService := payment.NewService(paymentRepo, sharedProducer, logger)
		paymentHandler := payment.NewHandler(paymentService, logger)
		paymentHandler.RegisterRoutes(protected)

		cardRepo := card.NewRepository(db, logger)
		cardService := card.NewService(cardRepo, sharedProducer, logger)
		cardHandler := card.NewHandler(cardService, logger)
		cardHandler.RegisterRoutes(protected)

		kycRepo := kyc.NewRepository(db, logger)
		kycService := kyc.NewService(kycRepo, sharedProducer, logger)
		kycHandler := kyc.NewHandler(kycService, logger)
		// Webhook must be on public (unauthenticated) group.
		kycHandler.RegisterWebhookRoute(v1)
		// Authenticated KYC routes.
		kycHandler.RegisterRoutes(protected)

		notifRepo := notification.NewRepository(db, logger)
		notifService := notification.NewService(notifRepo, sharedProducer, logger)
		notifHandler := notification.NewHandler(notifService, logger)
		notifHandler.RegisterRoutes(protected)

		cryptoRepo := crypto.NewRepository(db, logger)
		cryptoService := crypto.NewService(cryptoRepo, sharedProducer, logger)
		cryptoHandler := crypto.NewHandler(cryptoService, logger)
		cryptoHandler.RegisterRoutes(protected)
	}

	// Start HTTP server with graceful shutdown.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine.
	go func() {
		logger.Info("HTTP server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received, draining connections...")

	// Give outstanding requests 30 seconds to complete.
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	// Close resources.
	if db != nil {
		db.Close()
	}
	_ = redisClient.Close()

	logger.Info("server shut down gracefully")
}

// initJWTManager initializes the JWT manager, falling back to no-op
// if keys are not configured (development mode).
func initJWTManager(cfg *common.ServiceConfig, logger *zap.Logger) (*tpcrypto.JWTManager, error) {
	if cfg.JWTPrivateKeyPath == "" && cfg.JWTPublicKeyPath == "" {
		logger.Warn("JWT keys not configured, running without token signing")
		return nil, fmt.Errorf("no JWT keys configured")
	}
	return tpcrypto.NewJWTManager(cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath, cfg.JWTIssuer)
}

// healthCheck returns a handler that reports the service health status.
func healthCheck(db *database.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "healthy"
		checks := make(map[string]string)

		if db != nil {
			if err := db.HealthCheck(c.Request.Context()); err != nil {
				checks["database"] = "unhealthy: " + err.Error()
				status = "degraded"
			} else {
				checks["database"] = "healthy"
			}
		}

		if err := redis.Ping(c.Request.Context()).Err(); err != nil {
			checks["redis"] = "unhealthy: " + err.Error()
			status = "degraded"
		} else {
			checks["redis"] = "healthy"
		}

		httpStatus := http.StatusOK
		if status != "healthy" {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, gin.H{
			"status":  status,
			"service": "api-gateway",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"checks":  checks,
		})
	}
}

// readinessCheck verifies the service is ready to accept traffic.
func readinessCheck(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db != nil {
			if err := db.HealthCheck(c.Request.Context()); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "reason": "database unavailable"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"ready": true})
	}
}

