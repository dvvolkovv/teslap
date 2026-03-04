package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the bcrypt work factor for password hashing.
	// 12 provides good security without excessive latency (~250ms).
	BcryptCost = 12
)

// HashPassword hashes a password using bcrypt with the configured cost factor.
func HashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword checks a plaintext password against a bcrypt hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSalt creates a cryptographically random salt of the given byte length.
func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	return hex.EncodeToString(salt), nil
}

// SHA256Hash computes a SHA-256 hash of the input and returns the hex-encoded result.
func SHA256Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// HMACSHA256 computes an HMAC-SHA256 digest for webhook signature verification.
func HMACSHA256(message, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyHMAC compares a received HMAC digest against the expected value
// using constant-time comparison to prevent timing attacks.
func VerifyHMAC(message, secret []byte, receivedDigest string) bool {
	expected := HMACSHA256(message, secret)
	return hmac.Equal([]byte(expected), []byte(receivedDigest))
}

// GenerateSecureToken creates a cryptographically random token of the specified byte length,
// returned as a hex-encoded string. Used for refresh tokens and opaque session tokens.
func GenerateSecureToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
