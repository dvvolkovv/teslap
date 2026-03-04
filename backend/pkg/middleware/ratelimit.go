package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Rate limits per tier as defined in API contracts.
var tierLimits = map[string]int{
	"basic":           60,
	"standard":        100,
	"premium":         300,
	"unauthenticated": 20,
}

// RateLimiter provides Redis-backed sliding window rate limiting.
type RateLimiter struct {
	client *redis.Client
	logger *zap.Logger
	window time.Duration
}

// NewRateLimiter creates a rate limiter backed by Redis sorted sets.
func NewRateLimiter(client *redis.Client, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		client: client,
		logger: logger,
		window: 1 * time.Minute,
	}
}

// Middleware returns a Gin middleware that enforces per-user rate limits.
// The limit is determined by the user's account tier.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine the rate limit key and limit.
		var key string
		var limit int

		userID, exists := c.Get("user_id")
		if exists {
			key = fmt.Sprintf("ratelimit:user:%s", userID)
			tier, _ := c.Get("tier")
			tierStr, ok := tier.(string)
			if !ok || tierStr == "" {
				tierStr = "basic"
			}
			limit = tierLimits[tierStr]
			if limit == 0 {
				limit = tierLimits["basic"]
			}
		} else {
			key = fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
			limit = tierLimits["unauthenticated"]
		}

		now := time.Now()
		windowStart := now.Add(-rl.window)

		ctx := c.Request.Context()
		allowed, remaining, err := rl.checkRate(ctx, key, now, windowStart, limit)
		if err != nil {
			// If Redis is down, allow the request but log the failure.
			rl.logger.Warn("rate limit check failed, allowing request", zap.Error(err))
			c.Next()
			return
		}

		// Set rate limit headers as per API contract.
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now.Add(rl.window).Unix(), 10))

		if !allowed {
			problem := common.NewRateLimitError()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, problem)
			return
		}

		c.Next()
	}
}

// checkRate uses a Redis sorted set sliding window to check and record the request.
func (rl *RateLimiter) checkRate(ctx context.Context, key string, now, windowStart time.Time, limit int) (allowed bool, remaining int, err error) {
	pipe := rl.client.Pipeline()

	// Remove entries outside the current window.
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixMicro(), 10))

	// Count current requests in the window.
	countCmd := pipe.ZCard(ctx, key)

	// Add the current request.
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixMicro()),
		Member: fmt.Sprintf("%d-%d", now.UnixNano(), now.UnixMicro()),
	})

	// Set TTL on the key to auto-cleanup.
	pipe.Expire(ctx, key, rl.window+time.Second)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := int(countCmd.Val())
	remaining = limit - count - 1
	if remaining < 0 {
		remaining = 0
	}

	return count < limit, remaining, nil
}
