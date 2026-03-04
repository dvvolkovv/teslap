// Package main is the entry point for the TeslaPay Auth Service.
// It runs as a standalone microservice handling authentication,
// session management, and device binding.
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

	"github.com/teslapay/backend/internal/auth"
	"github.com/teslapay/backend/internal/common"
	tpcrypto "github.com/teslapay/backend/pkg/crypto"
	"github.com/teslapay/backend/pkg/database"
	"github.com/teslapay/backend/pkg/middleware"
)

func main() {
	cfg, err := common.LoadBaseConfig("auth-service")
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

	// Database connection.
	if cfg.DatabaseURL == "" {
		logger.Fatal("DATABASE_URL is required")
	}
	db, err := database.New(ctx, database.DefaultConfig(cfg.DatabaseURL), logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// JWT manager.
	jwtManager, err := tpcrypto.NewJWTManager(
		cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath, cfg.JWTIssuer,
	)
	if err != nil {
		logger.Fatal("failed to initialize JWT manager", zap.Error(err))
	}

	// Service setup.
	repo := auth.NewRepository(db, logger)
	service := auth.NewService(repo, jwtManager, logger, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	handler := auth.NewHandler(service, logger)

	// Router.
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLogger(logger))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth-service",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("auth-service listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down auth-service")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 15*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	logger.Info("auth-service stopped")
}
