package common

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// NewLogger creates a structured zap logger configured for the given environment.
// Production uses JSON encoding; development uses console encoding.
func NewLogger(serviceName, environment, level string) (*zap.Logger, error) {
	var cfg zap.Config
	if environment == "production" || environment == "staging" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		parsedLevel = zapcore.InfoLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(parsedLevel)

	logger, err := cfg.Build(
		zap.Fields(
			zap.String("service", serviceName),
			zap.String("environment", environment),
			zap.Int("pid", os.Getpid()),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// ContextWithLogger stores a logger instance in the context.
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext retrieves the logger from context. Returns a no-op logger if not found.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return zap.NewNop()
}

// ContextWithRequestID stores a request ID in the context.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext retrieves the request ID from context.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextWithUserID stores a user ID in the context.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext retrieves the user ID from context.
func UserIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}
