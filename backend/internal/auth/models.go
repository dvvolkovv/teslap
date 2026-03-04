// Package auth implements the authentication and session management domain
// for TeslaPay, including registration, login, JWT token management, and
// device binding.
package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserCredentials stores authentication credentials separately from the user
// profile for security isolation (auth_db vs account_db).
type UserCredentials struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	Email          string     `json:"email" db:"email"`
	EmailVerified  bool       `json:"email_verified" db:"email_verified"`
	Phone          string     `json:"phone" db:"phone"`
	PhoneVerified  bool       `json:"phone_verified" db:"phone_verified"`
	PasswordHash   string     `json:"-" db:"password_hash"`
	PasswordSalt   string     `json:"-" db:"password_salt"`
	FailedAttempts int        `json:"failed_attempts" db:"failed_attempts"`
	LockedUntil    *time.Time `json:"locked_until,omitempty" db:"locked_until"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	Status         string     `json:"status" db:"status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Session represents an active user session bound to a specific device.
type Session struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	DeviceID      uuid.UUID `json:"device_id" db:"device_id"`
	AccessTokenJTI string   `json:"access_token_jti" db:"access_token_jti"`
	IPAddress     string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent     string    `json:"user_agent,omitempty" db:"user_agent"`
	Location      string    `json:"location,omitempty" db:"location"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastActiveAt  time.Time `json:"last_active_at" db:"last_active_at"`
}

// SessionListItem is the public representation of a session for the sessions list API.
type SessionListItem struct {
	ID           string    `json:"id"`
	DeviceName   string    `json:"device_name"`
	DeviceType   string    `json:"device_type"`
	IPAddress    string    `json:"ip_address"`
	Location     string    `json:"location"`
	LastActiveAt time.Time `json:"last_active_at"`
	IsCurrent    bool      `json:"is_current"`
}

// Device represents a registered user device.
type Device struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	DeviceName       string     `json:"device_name" db:"device_name"`
	DeviceType       string     `json:"device_type" db:"device_type"`
	DeviceFingerprint string    `json:"device_fingerprint" db:"device_fingerprint"`
	PushToken        string     `json:"push_token,omitempty" db:"push_token"`
	BiometricKey     string     `json:"-" db:"biometric_key"`
	IsTrusted        bool       `json:"is_trusted" db:"is_trusted"`
	RegisteredAt     time.Time  `json:"registered_at" db:"registered_at"`
	LastSeenAt       *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
}

// RefreshToken is stored as a SHA-256 hash; the plaintext token is never persisted.
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	DeviceID  uuid.UUID  `json:"device_id" db:"device_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// -- Request/Response Types matching the API contracts --

// RegisterRequest matches POST /api/v1/auth/register.
type RegisterRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Phone    string         `json:"phone" binding:"required"`
	Password string         `json:"password" binding:"required,min=8"`
	Language string         `json:"language" binding:"required"`
	Device   DeviceInfo     `json:"device" binding:"required"`
	Consent  ConsentInfo    `json:"consent" binding:"required"`
}

// DeviceInfo holds device registration data.
type DeviceInfo struct {
	DeviceID   string `json:"device_id" binding:"required"`
	DeviceName string `json:"device_name" binding:"required"`
	DeviceType string `json:"device_type" binding:"required,oneof=ios android"`
	PushToken  string `json:"push_token,omitempty"`
}

// ConsentInfo holds the consent flags required during registration.
type ConsentInfo struct {
	TermsAccepted   bool `json:"terms_accepted" binding:"required"`
	PrivacyAccepted bool `json:"privacy_accepted" binding:"required"`
	MarketingOptIn  bool `json:"marketing_opt_in"`
}

// RegisterResponse matches the 201 response for registration.
type RegisterResponse struct {
	UserID                string `json:"user_id"`
	EmailVerificationSent bool   `json:"email_verification_sent"`
	PhoneVerificationSent bool   `json:"phone_verification_sent"`
}

// LoginRequest matches POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required"`
	Device   DeviceInfo `json:"device" binding:"required"`
}

// LoginResponse matches the 200 response for login.
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	User         LoginUser `json:"user"`
}

// LoginUser is the user info embedded in the login response.
type LoginUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	KYCStatus string `json:"kyc_status"`
	Tier      string `json:"tier"`
}

// RefreshRequest matches POST /api/v1/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse matches the 200 response for token refresh.
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// LogoutRequest matches POST /api/v1/auth/logout.
type LogoutRequest struct {
	AllSessions bool `json:"all_sessions"`
}

// VerifyRequest matches POST /api/v1/auth/verify-email and verify-phone.
type VerifyRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required,len=6"`
}

// VerifyResponse is the response for email/phone verification.
type VerifyResponse struct {
	EmailVerified bool `json:"email_verified,omitempty"`
	PhoneVerified bool `json:"phone_verified,omitempty"`
}

// Account limits for lockout.
const (
	MaxFailedAttempts = 5
	LockoutDuration   = 30 * time.Minute
)
