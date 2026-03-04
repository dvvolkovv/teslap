package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Recovery is a Gin middleware that recovers from panics and returns a
// structured RFC 7807 error response instead of crashing the server.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				requestID := common.RequestIDFromContext(c.Request.Context())
				logger.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("request_id", requestID),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				problem := common.NewInternalError(requestID)
				c.AbortWithStatusJSON(http.StatusInternalServerError, problem)
			}
		}()

		c.Next()
	}
}
