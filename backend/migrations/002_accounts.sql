-- Migration 002: Account management tables
-- Tables: account_tiers, users, accounts, sub_accounts, beneficiaries
-- Derived from internal/account/models.go

-- ---------------------------------------------------------------------------
-- account_tiers: feature limit definitions per tier level
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS account_tiers (
    id                    UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    name                  TEXT           NOT NULL,
    monthly_fee           NUMERIC(18,4)  NOT NULL DEFAULT 0,
    daily_transfer_limit  NUMERIC(18,4),
    monthly_transfer_limit NUMERIC(18,4),
    daily_card_limit      NUMERIC(18,4),
    monthly_atm_limit     NUMERIC(18,4),
    free_atm_withdrawals  INT            NOT NULL DEFAULT 0,
    fx_markup_percent     NUMERIC(5,4)   NOT NULL DEFAULT 0,
    max_sub_accounts      INT            NOT NULL DEFAULT 5,
    created_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_account_tiers_name UNIQUE (name)
);

-- Seed default tiers — idempotent via ON CONFLICT DO NOTHING.
INSERT INTO account_tiers (
    id, name, monthly_fee,
    daily_transfer_limit, monthly_transfer_limit,
    daily_card_limit, monthly_atm_limit,
    free_atm_withdrawals, fx_markup_percent, max_sub_accounts
) VALUES
    (gen_random_uuid(), 'Standard', 0,     5000,   20000,   2000,  500,   2,  0.015, 3),
    (gen_random_uuid(), 'Premium',  9.99,  50000,  200000,  10000, 2000,  10, 0.005, 10),
    (gen_random_uuid(), 'Business', 29.99, 500000, 2000000, 50000, 10000, 0,  0.010, 25)
ON CONFLICT (name) DO NOTHING;

-- ---------------------------------------------------------------------------
-- users: user profile data (separate from auth credentials)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id   TEXT,
    tier_id       UUID        REFERENCES account_tiers(id),
    first_name    TEXT,
    last_name     TEXT,
    date_of_birth DATE,
    nationality   TEXT,
    tax_residency TEXT,
    address_line1 TEXT,
    address_line2 TEXT,
    city          TEXT,
    postal_code   TEXT,
    country       TEXT        NOT NULL,
    language      TEXT        NOT NULL DEFAULT 'en',
    kyc_status    TEXT        NOT NULL DEFAULT 'none',
    kyc_level     INT         NOT NULL DEFAULT 0,
    risk_score    INT         NOT NULL DEFAULT 0,
    status        TEXT        NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT uq_users_external_id UNIQUE (external_id)
);

CREATE INDEX IF NOT EXISTS idx_users_tier_id    ON users(tier_id);
CREATE INDEX IF NOT EXISTS idx_users_kyc_status ON users(kyc_status);
CREATE INDEX IF NOT EXISTS idx_users_status     ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;

-- ---------------------------------------------------------------------------
-- accounts: top-level account container (one per user)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS accounts (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_number TEXT        NOT NULL,
    status         TEXT        NOT NULL DEFAULT 'active',
    opened_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at      TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_accounts_account_number UNIQUE (account_number)
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_status  ON accounts(status);

-- ---------------------------------------------------------------------------
-- sub_accounts: single-currency account with optional IBAN (LT format)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS sub_accounts (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id       UUID        NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    currency         CHAR(3)     NOT NULL,
    iban             TEXT,
    bic              TEXT,
    ledger_account_id UUID,
    status           TEXT        NOT NULL DEFAULT 'active',
    is_default       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_sub_accounts_iban           UNIQUE (iban),
    CONSTRAINT uq_sub_accounts_account_currency UNIQUE (account_id, currency)
);

CREATE INDEX IF NOT EXISTS idx_sub_accounts_account_id ON sub_accounts(account_id);
CREATE INDEX IF NOT EXISTS idx_sub_accounts_currency   ON sub_accounts(currency);

-- ---------------------------------------------------------------------------
-- beneficiaries: saved payees / recipients for payments
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS beneficiaries (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name              TEXT        NOT NULL,
    iban              TEXT        NOT NULL,
    bic               TEXT,
    bank_name         TEXT,
    default_reference TEXT,
    is_internal       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_beneficiaries_user_id    ON beneficiaries(user_id);
CREATE INDEX IF NOT EXISTS idx_beneficiaries_iban       ON beneficiaries(iban);
CREATE INDEX IF NOT EXISTS idx_beneficiaries_deleted_at ON beneficiaries(deleted_at) WHERE deleted_at IS NOT NULL;
