// Package main is the entry point for the TeslaPay Card Service.
// Handles virtual/physical card issuance, lifecycle management,
// spending controls, and authorization processing.
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

	"github.com/teslapay/backend/internal/card"
	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/pkg/database"
	"github.com/teslapay/backend/pkg/events"
	"github.com/teslapay/backend/pkg/middleware"
)

func main() {
	cfg, err := common.LoadBaseConfig("card-service")
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

	// Wire up card service components.
	cardRepo := card.NewRepository(db, logger)
	cardService := card.NewService(cardRepo, producer, logger)
	cardHandler := card.NewHandler(cardService, logger)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLogger(logger))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "card-service",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	v1 := router.Group("/api/v1")
	cardHandler.RegisterRoutes(v1)

	// Internal webhook endpoints for card processor (Enfuce).
	// These are called directly by the card processor, not via the API gateway.
	router.POST("/internal/webhooks/enfuce/authorization", notImplemented)
	router.POST("/internal/webhooks/enfuce/settlement", notImplemented)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,  // Card auth is latency-critical.
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("card-service listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down card-service")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}

// notImplemented returns 501 for webhook endpoints pending card processor integration.
func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"type":       "https://api.teslapay.eu/errors/not-implemented",
		"title":      "Not Implemented",
		"status":     501,
		"detail":     "This webhook endpoint is pending card processor integration",
		"error_code": "NOT_IMPLEMENTED",
	})
}
