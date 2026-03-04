package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// RequestLogger is a Gin middleware that logs every HTTP request with
// structured fields including timing, status, and correlation IDs.
// It injects a request-scoped logger into the context for downstream use.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Use the client-provided request ID or generate one.
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Echo request ID back in the response.
		c.Header("X-Request-ID", requestID)

		// Inject request ID and logger into context.
		ctx := common.ContextWithRequestID(c.Request.Context(), requestID)
		reqLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)
		ctx = common.ContextWithLogger(ctx, reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// Process the request.
		c.Next()

		// Log the completed request.
		duration := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.Int("response_size", c.Writer.Size()),
		}

		// Include user_id if authenticated.
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		if status >= 500 {
			reqLogger.Error("request completed with server error", fields...)
		} else if status >= 400 {
			reqLogger.Warn("request completed with client error", fields...)
		} else {
			reqLogger.Info("request completed", fields...)
		}
	}
}
