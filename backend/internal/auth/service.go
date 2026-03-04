package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	tpcrypto "github.com/teslapay/backend/pkg/crypto"
)

// Service implements authentication business logic.
type Service struct {
	repo       *Repository
	jwt        *tpcrypto.JWTManager
	logger     *zap.Logger
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewService creates a new auth service.
func NewService(
	repo *Repository,
	jwt *tpcrypto.JWTManager,
	logger *zap.Logger,
	accessTTL, refreshTTL time.Duration,
) *Service {
	return &Service{
		repo:       repo,
		jwt:        jwt,
		logger:     logger,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// Register creates a new user credential record and associated device.
// It returns a user ID and triggers email/phone verification (placeholder).
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Validate consent.
	if !req.Consent.TermsAccepted || !req.Consent.PrivacyAccepted {
		return nil, common.NewValidationError(
			"Terms and privacy policy must be accepted", "/api/v1/auth/register",
		)
	}

	// Check if email or phone already exists.
	existing, err := s.repo.GetCredentialsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check existing email: %w", err)
	}
	if existing != nil {
		return nil, common.NewConflictError("An account with this email already exists")
	}

	existingPhone, err := s.repo.GetCredentialsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("check existing phone: %w", err)
	}
	if existingPhone != nil {
		return nil, common.NewConflictError("An account with this phone number already exists")
	}

	// Hash the password.
	passwordHash, err := tpcrypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	salt, err := tpcrypto.GenerateSalt(16)
	if err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	userID := uuid.New()
	now := time.Now().UTC()

	creds := &UserCredentials{
		ID:           uuid.New(),
		UserID:       userID,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: passwordHash,
		PasswordSalt: salt,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateCredentials(ctx, creds); err != nil {
		return nil, fmt.Errorf("create credentials: %w", err)
	}

	// Register the device.
	deviceUUID, err := uuid.Parse(req.Device.DeviceID)
	if err != nil {
		deviceUUID = uuid.New()
	}
	device := &Device{
		ID:               deviceUUID,
		UserID:           userID,
		DeviceName:       req.Device.DeviceName,
		DeviceType:       req.Device.DeviceType,
		DeviceFingerprint: req.Device.DeviceID,
		PushToken:        req.Device.PushToken,
		IsTrusted:        true, // First device is auto-trusted.
		RegisteredAt:     now,
	}
	if err := s.repo.CreateDevice(ctx, device); err != nil {
		return nil, fmt.Errorf("register device: %w", err)
	}

	s.logger.Info("user registered",
		zap.String("user_id", userID.String()),
		zap.String("email", req.Email),
	)

	// In production: trigger email and SMS OTP verification here
	// via the notification service Kafka topic.

	return &RegisterResponse{
		UserID:                userID.String(),
		EmailVerificationSent: true,
		PhoneVerificationSent: true,
	}, nil
}

// Login authenticates a user with email/password and returns JWT tokens.
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	creds, err := s.repo.GetCredentialsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("lookup credentials: %w", err)
	}
	if creds == nil {
		return nil, common.NewAuthError(common.ErrCodeAuthMissing, "Invalid email or password")
	}

	// Check if account is locked.
	if creds.Status == "locked" || creds.Status == "suspended" {
		return nil, common.NewAuthError(common.ErrCodeAuthForbidden,
			"Account is locked. Please contact support.")
	}
	if creds.LockedUntil != nil && time.Now().UTC().Before(*creds.LockedUntil) {
		return nil, common.NewAuthError(common.ErrCodeAuthForbidden,
			"Account is temporarily locked due to too many failed attempts")
	}

	// Verify password.
	if !tpcrypto.VerifyPassword(req.Password, creds.PasswordHash) {
		// Increment failed attempts.
		newFailedCount := creds.FailedAttempts + 1
		var lockUntil *time.Time
		if newFailedCount >= MaxFailedAttempts {
			t := time.Now().UTC().Add(LockoutDuration)
			lockUntil = &t
		}
		_ = s.repo.UpdateFailedAttempts(ctx, creds.ID, newFailedCount, lockUntil)

		s.logger.Warn("failed login attempt",
			zap.String("email", req.Email),
			zap.Int("attempt", newFailedCount),
		)

		return nil, common.NewAuthError(common.ErrCodeAuthMissing, "Invalid email or password")
	}

	// Reset failed attempts on successful login.
	_ = s.repo.UpdateFailedAttempts(ctx, creds.ID, 0, nil)

	// Check if this device is trusted.
	device, err := s.repo.GetDeviceByFingerprint(ctx, creds.UserID, req.Device.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("lookup device: %w", err)
	}
	if device == nil {
		// New device detected -- in production, trigger OTP challenge.
		s.logger.Info("new device detected during login",
			zap.String("user_id", creds.UserID.String()),
			zap.String("device_id", req.Device.DeviceID),
		)
		// For MVP, auto-trust and register the new device.
		deviceUUID, _ := uuid.Parse(req.Device.DeviceID)
		if deviceUUID == uuid.Nil {
			deviceUUID = uuid.New()
		}
		device = &Device{
			ID:                deviceUUID,
			UserID:            creds.UserID,
			DeviceName:        req.Device.DeviceName,
			DeviceType:        req.Device.DeviceType,
			DeviceFingerprint: req.Device.DeviceID,
			IsTrusted:         true,
			RegisteredAt:      time.Now().UTC(),
		}
		_ = s.repo.CreateDevice(ctx, device)
	}

	// Generate tokens.
	sessionID := uuid.New()

	accessToken, err := s.jwt.GenerateAccessToken(
		creds.UserID.String(), creds.Email,
		"standard", // Tier fetched from account service in production.
		"pending",  // KYC status fetched from KYC service in production.
		device.ID.String(),
		sessionID.String(),
		s.accessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	// Generate opaque refresh token.
	refreshTokenPlain, err := tpcrypto.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshTokenHash := tpcrypto.SHA256Hash(refreshTokenPlain)

	now := time.Now().UTC()
	rt := &RefreshToken{
		ID:        uuid.New(),
		UserID:    creds.UserID,
		DeviceID:  device.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	// Create session record.
	jti, _ := s.jwt.GetJTI(accessToken)
	session := &Session{
		ID:             sessionID,
		UserID:         creds.UserID,
		DeviceID:       device.ID,
		AccessTokenJTI: jti,
		ExpiresAt:      now.Add(s.refreshTTL),
		CreatedAt:      now,
		LastActiveAt:   now,
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	// Update last login timestamp.
	_ = s.repo.UpdateLastLogin(ctx, creds.ID, now)

	s.logger.Info("user logged in",
		zap.String("user_id", creds.UserID.String()),
		zap.String("session_id", sessionID.String()),
	)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenPlain,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTTL.Seconds()),
		User: LoginUser{
			ID:        creds.UserID.String(),
			Email:     creds.Email,
			FirstName: "", // Fetched from account service in production.
			KYCStatus: "pending",
			Tier:      "standard",
		},
	}, nil
}

// RefreshToken exchanges a valid refresh token for new access and refresh tokens.
// The old refresh token is revoked (rotation).
func (s *Service) RefreshToken(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error) {
	tokenHash := tpcrypto.SHA256Hash(req.RefreshToken)

	rt, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("lookup refresh token: %w", err)
	}
	if rt == nil {
		return nil, common.NewAuthError(common.ErrCodeAuthExpired, "Refresh token is invalid")
	}
	if rt.RevokedAt != nil {
		s.logger.Warn("attempted reuse of revoked refresh token",
			zap.String("user_id", rt.UserID.String()),
			zap.String("token_id", rt.ID.String()),
		)
		// Potential token theft -- revoke all sessions for this user.
		_ = s.repo.RevokeAllRefreshTokens(ctx, rt.UserID)
		return nil, common.NewAuthError(common.ErrCodeAuthExpired,
			"Refresh token has been revoked. All sessions terminated for security.")
	}
	if time.Now().UTC().After(rt.ExpiresAt) {
		return nil, common.NewAuthError(common.ErrCodeAuthExpired, "Refresh token has expired")
	}

	// Revoke the old refresh token.
	if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
		return nil, fmt.Errorf("revoke old refresh token: %w", err)
	}

	// Load credentials for token claims.
	creds, err := s.repo.GetCredentialsByUserID(ctx, rt.UserID)
	if err != nil || creds == nil {
		return nil, common.NewAuthError(common.ErrCodeAuthMissing, "User not found")
	}

	// Generate new tokens.
	newAccessToken, err := s.jwt.GenerateAccessToken(
		creds.UserID.String(), creds.Email,
		"standard", "pending",
		rt.DeviceID.String(), uuid.New().String(),
		s.accessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("generate new access token: %w", err)
	}

	newRefreshPlain, err := tpcrypto.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate new refresh token: %w", err)
	}
	newRefreshHash := tpcrypto.SHA256Hash(newRefreshPlain)

	now := time.Now().UTC()
	newRT := &RefreshToken{
		ID:        uuid.New(),
		UserID:    rt.UserID,
		DeviceID:  rt.DeviceID,
		TokenHash: newRefreshHash,
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.repo.CreateRefreshToken(ctx, newRT); err != nil {
		return nil, fmt.Errorf("store new refresh token: %w", err)
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshPlain,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTTL.Seconds()),
	}, nil
}

// ValidateToken validates a JWT access token and returns the claims.
// Used by other services for token introspection.
func (s *Service) ValidateToken(tokenString string) (*tpcrypto.TokenClaims, error) {
	return s.jwt.ValidateAccessToken(tokenString)
}
