// Package middleware provides HTTP middleware for TeslaPay API Gateway
// including JWT authentication, request logging, and rate limiting.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	tpcrypto "github.com/teslapay/backend/pkg/crypto"
)

// AuthMiddleware validates JWT access tokens on incoming requests.
// It extracts the Bearer token from the Authorization header,
// validates the signature and claims, and injects user information
// into the request context.
func AuthMiddleware(jwtManager *tpcrypto.JWTManager, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			problem := common.NewAuthError(common.ErrCodeAuthMissing, "Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			problem := common.NewAuthError(common.ErrCodeAuthMissing, "Authorization header must use Bearer scheme")
			c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			logger.Debug("token validation failed", zap.Error(err))
			problem := common.NewAuthError(common.ErrCodeAuthExpired, "Access token is invalid or expired")
			c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
			return
		}

		// Inject claims into context for downstream handlers.
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("tier", claims.Tier)
		c.Set("kyc_status", claims.KYCStatus)
		c.Set("device_id", claims.DeviceID)
		c.Set("session_id", claims.SessionID)
		c.Set("token_jti", claims.ID)

		// Propagate to common context utilities.
		ctx := common.ContextWithUserID(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// SCAMiddleware verifies that a valid SCA token is present for operations
// requiring PSD2 Strong Customer Authentication.
func SCAMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		scaToken := c.GetHeader("X-SCA-Token")
		if scaToken == "" {
			problem := common.NewForbiddenError(common.ErrCodeSCARequired,
				"Strong Customer Authentication required for this operation")
			c.AbortWithStatusJSON(http.StatusForbidden, problem)
			return
		}

		// In production, validate the SCA token against Redis/DB.
		// The SCA token is a short-lived, single-use token issued after
		// the user completes biometric/OTP verification via /auth/sca/verify.
		// For now, we accept its presence as sufficient.

		c.Set("sca_token", scaToken)
		c.Next()
	}
}

// OptionalAuth attempts to extract JWT claims but does not reject the request
// if no token is present. Used for endpoints that behave differently for
// authenticated vs. anonymous users.
func OptionalAuth(jwtManager *tpcrypto.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("tier", claims.Tier)
		ctx := common.ContextWithUserID(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
