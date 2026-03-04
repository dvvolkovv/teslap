// Package main is the entry point for the TeslaPay Ledger Service.
// The ledger is the financial core of the platform, handling double-entry
// bookkeeping, balance management, and event sourcing.
//
// This service is INTERNAL ONLY -- it is never directly exposed to external clients.
// All access is via gRPC from other services (Payment, Card, Crypto).
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
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/internal/ledger"
	"github.com/teslapay/backend/pkg/database"
	"github.com/teslapay/backend/pkg/events"
	"github.com/teslapay/backend/pkg/middleware"
)

func main() {
	cfg, err := common.LoadBaseConfig("ledger-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger, err := common.NewLogger(cfg.ServiceName, cfg.Environment, cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Database (required).
	if cfg.DatabaseURL == "" {
		logger.Fatal("DATABASE_URL is required for ledger service")
	}
	dbCfg := database.DefaultConfig(cfg.DatabaseURL)
	dbCfg.MaxConns = 50 // Ledger has highest write volume.
	db, err := database.New(ctx, dbCfg, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Kafka producer for ledger events.
	var producer *events.Producer
	if len(cfg.KafkaBrokers) > 0 {
		producerCfg := events.DefaultProducerConfig(cfg.KafkaBrokers)
		producer = events.NewProducer(producerCfg, logger)
		defer func() { _ = producer.Close() }()
	}

	// Service setup.
	repo := ledger.NewRepository(db, logger)
	service := ledger.NewService(repo, producer, logger)

	// HTTP server for health checks and internal REST API.
	// In production, the primary interface is gRPC. The HTTP server
	// provides health/readiness endpoints and optional REST fallback.
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLogger(logger))

	router.GET("/health", func(c *gin.Context) {
		dbErr := db.HealthCheck(c.Request.Context())
		status := "healthy"
		if dbErr != nil {
			status = "unhealthy"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"service":  "ledger-service",
			"time":     time.Now().UTC().Format(time.RFC3339),
			"database": status,
		})
	})

	// Internal API for direct HTTP access (dev/testing).
	internal := router.Group("/internal/ledger")
	{
		internal.POST("/post", func(c *gin.Context) {
			var req ledger.PostTransactionRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			je, err := service.PostTransaction(c.Request.Context(), &req)
			if err != nil {
				if problem, ok := err.(*common.ProblemDetail); ok {
					c.JSON(problem.Status, problem)
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, je)
		})

		internal.GET("/balance/:account_id", func(c *gin.Context) {
			accountID, err := parseUUID(c.Param("account_id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
				return
			}
			bal, err := service.GetBalance(c.Request.Context(), accountID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, bal)
		})
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("ledger-service listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// TODO: Start gRPC server on a separate port for inter-service communication.

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down ledger-service")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 15*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	logger.Info("ledger-service stopped")
}

func parseUUID(s string) ([16]byte, error) {
	if len(s) != 36 {
		return [16]byte{}, fmt.Errorf("invalid UUID: %s", s)
	}
	var result [16]byte
	hexChars := make([]byte, 0, 32)
	for _, c := range s {
		if c != '-' {
			hexChars = append(hexChars, byte(c))
		}
	}
	if len(hexChars) != 32 {
		return result, fmt.Errorf("invalid UUID: %s", s)
	}
	for i := 0; i < 16; i++ {
		result[i] = hexByte(hexChars[i*2])<<4 | hexByte(hexChars[i*2+1])
	}
	return result, nil
}

func hexByte(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

// Ensure the service variable is used to avoid compile error.
var _ = (*ledger.Service)(nil)
