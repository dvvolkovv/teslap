// Package main is the entry point for the TeslaPay Notification Service.
// Handles push notifications, email, SMS, and in-app notifications.
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
	"github.com/teslapay/backend/internal/notification"
	"github.com/teslapay/backend/pkg/database"
	"github.com/teslapay/backend/pkg/events"
	"github.com/teslapay/backend/pkg/middleware"
)

func main() {
	cfg, err := common.LoadBaseConfig("notification-service")
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

	if cfg.DatabaseURL == "" {
		logger.Fatal("DATABASE_URL is required")
	}
	db, err := database.New(ctx, database.DefaultConfig(cfg.DatabaseURL), logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Kafka producer (optional — gracefully degrade if unavailable).
	var producer *events.Producer
	if len(cfg.KafkaBrokers) > 0 {
		producerCfg := events.DefaultProducerConfig(cfg.KafkaBrokers)
		producer = events.NewProducer(producerCfg, logger)
		defer func() {
			if err := producer.Close(); err != nil {
				logger.Error("failed to close kafka producer", zap.Error(err))
			}
		}()
		logger.Info("kafka producer initialized", zap.Strings("brokers", cfg.KafkaBrokers))
	} else {
		logger.Warn("KAFKA_BROKERS not set, events will not be published")
	}

	// Wire up notification service components.
	notifRepo := notification.NewRepository(db, logger)
	notifService := notification.NewService(notifRepo, producer, logger)
	notifHandler := notification.NewHandler(notifService, logger)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLogger(logger))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "notification-service",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	v1 := router.Group("/api/v1")
	notifHandler.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info("notification-service listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down notification-service")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 15*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}
