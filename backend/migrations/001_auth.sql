-- Migration 001: Authentication tables
-- Tables: user_credentials, devices, sessions, refresh_tokens
-- Derived from internal/auth/models.go

-- ---------------------------------------------------------------------------
-- user_credentials: stores authentication identity, separate from user profile
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_credentials (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL,
    email            TEXT        NOT NULL,
    email_verified   BOOLEAN     NOT NULL DEFAULT FALSE,
    phone            TEXT,
    phone_verified   BOOLEAN     NOT NULL DEFAULT FALSE,
    password_hash    TEXT,
    password_salt    TEXT,
    failed_attempts  INT         NOT NULL DEFAULT 0,
    locked_until     TIMESTAMPTZ,
    last_login_at    TIMESTAMPTZ,
    status           TEXT        NOT NULL DEFAULT 'pending',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_credentials_email UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_user_credentials_user_id ON user_credentials(user_id);
CREATE INDEX IF NOT EXISTS idx_user_credentials_email   ON user_credentials(email);
CREATE INDEX IF NOT EXISTS idx_user_credentials_status  ON user_credentials(status);

-- ---------------------------------------------------------------------------
-- devices: registered user devices (iOS / Android), bound to a user
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS devices (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID        NOT NULL REFERENCES user_credentials(user_id) ON DELETE CASCADE,
    device_name        TEXT,
    device_type        TEXT,
    device_fingerprint TEXT,
    push_token         TEXT,
    biometric_key      TEXT,
    is_trusted         BOOLEAN     NOT NULL DEFAULT FALSE,
    registered_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at       TIMESTAMPTZ,
    CONSTRAINT uq_devices_fingerprint UNIQUE (device_fingerprint)
);

CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);

-- ---------------------------------------------------------------------------
-- sessions: active user sessions bound to a device and a JWT JTI
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS sessions (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL,
    device_id        UUID        NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    access_token_jti TEXT,
    ip_address       TEXT,
    user_agent       TEXT,
    location         TEXT,
    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id   ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_device_id ON sessions(device_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires   ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_jti       ON sessions(access_token_jti);

-- ---------------------------------------------------------------------------
-- refresh_tokens: hashed refresh tokens (plaintext is never stored)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    device_id  UUID        NOT NULL,
    token_hash TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id    ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_device_id  ON refresh_tokens(device_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires    ON refresh_tokens(expires_at);
