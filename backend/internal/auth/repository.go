package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the auth service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new auth repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreateCredentials inserts a new user credentials record.
func (r *Repository) CreateCredentials(ctx context.Context, creds *UserCredentials) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO user_credentials
			(id, user_id, email, email_verified, phone, phone_verified,
			 password_hash, password_salt, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		creds.ID, creds.UserID, creds.Email, creds.EmailVerified,
		creds.Phone, creds.PhoneVerified, creds.PasswordHash,
		creds.PasswordSalt, creds.Status, creds.CreatedAt, creds.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert credentials: %w", err)
	}
	return nil
}

// GetCredentialsByEmail looks up credentials by email address.
func (r *Repository) GetCredentialsByEmail(ctx context.Context, email string) (*UserCredentials, error) {
	var creds UserCredentials
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, email, email_verified, phone, phone_verified,
		       password_hash, password_salt, failed_attempts, locked_until,
		       last_login_at, status, created_at, updated_at
		FROM user_credentials
		WHERE email = $1
	`, email).Scan(
		&creds.ID, &creds.UserID, &creds.Email, &creds.EmailVerified,
		&creds.Phone, &creds.PhoneVerified, &creds.PasswordHash,
		&creds.PasswordSalt, &creds.FailedAttempts, &creds.LockedUntil,
		&creds.LastLoginAt, &creds.Status, &creds.CreatedAt, &creds.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query credentials by email: %w", err)
	}
	return &creds, nil
}

// GetCredentialsByPhone looks up credentials by phone number.
func (r *Repository) GetCredentialsByPhone(ctx context.Context, phone string) (*UserCredentials, error) {
	var creds UserCredentials
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, email, email_verified, phone, phone_verified,
		       password_hash, password_salt, failed_attempts, locked_until,
		       last_login_at, status, created_at, updated_at
		FROM user_credentials
		WHERE phone = $1
	`, phone).Scan(
		&creds.ID, &creds.UserID, &creds.Email, &creds.EmailVerified,
		&creds.Phone, &creds.PhoneVerified, &creds.PasswordHash,
		&creds.PasswordSalt, &creds.FailedAttempts, &creds.LockedUntil,
		&creds.LastLoginAt, &creds.Status, &creds.CreatedAt, &creds.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query credentials by phone: %w", err)
	}
	return &creds, nil
}

// GetCredentialsByUserID looks up credentials by user UUID.
func (r *Repository) GetCredentialsByUserID(ctx context.Context, userID uuid.UUID) (*UserCredentials, error) {
	var creds UserCredentials
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, email, email_verified, phone, phone_verified,
		       password_hash, password_salt, failed_attempts, locked_until,
		       last_login_at, status, created_at, updated_at
		FROM user_credentials
		WHERE user_id = $1
	`, userID).Scan(
		&creds.ID, &creds.UserID, &creds.Email, &creds.EmailVerified,
		&creds.Phone, &creds.PhoneVerified, &creds.PasswordHash,
		&creds.PasswordSalt, &creds.FailedAttempts, &creds.LockedUntil,
		&creds.LastLoginAt, &creds.Status, &creds.CreatedAt, &creds.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query credentials by user_id: %w", err)
	}
	return &creds, nil
}

// UpdateFailedAttempts updates the failed login attempt count and optional lockout time.
func (r *Repository) UpdateFailedAttempts(ctx context.Context, credID uuid.UUID, count int, lockedUntil *time.Time) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE user_credentials
		SET failed_attempts = $2, locked_until = $3, updated_at = NOW()
		WHERE id = $1
	`, credID, count, lockedUntil)
	return err
}

// UpdateLastLogin updates the last login timestamp.
func (r *Repository) UpdateLastLogin(ctx context.Context, credID uuid.UUID, at time.Time) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE user_credentials
		SET last_login_at = $2, updated_at = NOW()
		WHERE id = $1
	`, credID, at)
	return err
}

// CreateDevice registers a new device for a user.
func (r *Repository) CreateDevice(ctx context.Context, device *Device) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO devices
			(id, user_id, device_name, device_type, device_fingerprint,
			 push_token, biometric_key, is_trusted, registered_at, last_seen_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO NOTHING
	`,
		device.ID, device.UserID, device.DeviceName, device.DeviceType,
		device.DeviceFingerprint, device.PushToken, device.BiometricKey,
		device.IsTrusted, device.RegisteredAt, device.LastSeenAt,
	)
	if err != nil {
		return fmt.Errorf("insert device: %w", err)
	}
	return nil
}

// GetDeviceByFingerprint retrieves a device by user ID and fingerprint.
func (r *Repository) GetDeviceByFingerprint(ctx context.Context, userID uuid.UUID, fingerprint string) (*Device, error) {
	var d Device
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, device_name, device_type, device_fingerprint,
		       push_token, biometric_key, is_trusted, registered_at, last_seen_at
		FROM devices
		WHERE user_id = $1 AND device_fingerprint = $2
	`, userID, fingerprint).Scan(
		&d.ID, &d.UserID, &d.DeviceName, &d.DeviceType, &d.DeviceFingerprint,
		&d.PushToken, &d.BiometricKey, &d.IsTrusted, &d.RegisteredAt, &d.LastSeenAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query device: %w", err)
	}
	return &d, nil
}

// CreateSession inserts a new session record.
func (r *Repository) CreateSession(ctx context.Context, session *Session) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO sessions
			(id, user_id, device_id, access_token_jti, ip_address,
			 user_agent, location, expires_at, created_at, last_active_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		session.ID, session.UserID, session.DeviceID, session.AccessTokenJTI,
		session.IPAddress, session.UserAgent, session.Location,
		session.ExpiresAt, session.CreatedAt, session.LastActiveAt,
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

// CreateRefreshToken stores a hashed refresh token.
func (r *Repository) CreateRefreshToken(ctx context.Context, rt *RefreshToken) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, device_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, rt.ID, rt.UserID, rt.DeviceID, rt.TokenHash, rt.ExpiresAt, rt.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// GetRefreshTokenByHash retrieves a refresh token by its SHA-256 hash.
func (r *Repository) GetRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, device_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, hash).Scan(
		&rt.ID, &rt.UserID, &rt.DeviceID, &rt.TokenHash,
		&rt.ExpiresAt, &rt.RevokedAt, &rt.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query refresh token: %w", err)
	}
	return &rt, nil
}

// RevokeRefreshToken marks a refresh token as revoked.
func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1
	`, tokenID)
	return err
}

// RevokeAllRefreshTokens revokes all refresh tokens for a user (logout everywhere).
func (r *Repository) RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	return err
}

// DeleteSession removes a session record.
func (r *Repository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionID)
	return err
}

// DeleteAllSessions removes all sessions for a user.
func (r *Repository) DeleteAllSessions(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}
