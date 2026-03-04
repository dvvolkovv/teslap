// Package crypto provides cryptographic utilities for TeslaPay including
// JWT token management with RS256, password hashing, and HMAC verification.
package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles RS256 JWT token creation and validation.
type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

// TokenClaims represents the JWT claims for TeslaPay access tokens.
type TokenClaims struct {
	jwt.RegisteredClaims
	UserID    string `json:"user_id"`
	Email     string `json:"email,omitempty"`
	Tier      string `json:"tier,omitempty"`
	KYCStatus string `json:"kyc_status,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// NewJWTManager creates a JWTManager from PEM-encoded key files.
// Both keys must be RSA keys for RS256 signing.
func NewJWTManager(privateKeyPath, publicKeyPath, issuer string) (*JWTManager, error) {
	manager := &JWTManager{issuer: issuer}

	if privateKeyPath != "" {
		privBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read private key: %w", err)
		}
		block, _ := pem.Decode(privBytes)
		if block == nil {
			return nil, errors.New("failed to decode private key PEM")
		}
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// Try PKCS8 format.
			key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err2 != nil {
				return nil, fmt.Errorf("parse private key: %w (pkcs1: %v)", err2, err)
			}
			var ok bool
			privKey, ok = key.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("private key is not RSA")
			}
		}
		manager.privateKey = privKey
	}

	if publicKeyPath != "" {
		pubBytes, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read public key: %w", err)
		}
		block, _ := pem.Decode(pubBytes)
		if block == nil {
			return nil, errors.New("failed to decode public key PEM")
		}
		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse public key: %w", err)
		}
		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not RSA")
		}
		manager.publicKey = rsaPubKey
	}

	// If we have a private key but no public key, derive it.
	if manager.privateKey != nil && manager.publicKey == nil {
		manager.publicKey = &manager.privateKey.PublicKey
	}

	return manager, nil
}

// GenerateAccessToken creates a signed JWT access token with the given claims.
func (m *JWTManager) GenerateAccessToken(userID, email, tier, kycStatus, deviceID, sessionID string, ttl time.Duration) (string, error) {
	if m.privateKey == nil {
		return "", errors.New("private key not configured")
	}

	now := time.Now().UTC()
	jti := uuid.New().String()

	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
		UserID:    userID,
		Email:     email,
		Tier:      tier,
		KYCStatus: kycStatus,
		DeviceID:  deviceID,
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

// ValidateAccessToken parses and validates a JWT access token, returning its claims.
func (m *JWTManager) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	if m.publicKey == nil {
		return nil, errors.New("public key not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	}, jwt.WithIssuer(m.issuer))

	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GetJTI extracts the JTI (JWT ID) from a token string without full validation.
// Used for revocation checks against Redis.
func (m *JWTManager) GetJTI(tokenString string) (string, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenString, &TokenClaims{})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}
	return claims.ID, nil
}
