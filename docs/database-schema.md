# TeslaPay Database Schema Design

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Principal Software Architect, Dream Team

---

## 1. Design Principles

| Principle | Application |
|-----------|-------------|
| Database-per-Service | Each microservice owns its database; no cross-database joins |
| Immutable Ledger | Ledger entries are append-only; corrections via reversal entries |
| Financial Precision | All monetary amounts stored as `NUMERIC(19,4)` -- 15 integer digits + 4 decimal places |
| Event Sourcing | All financial state changes stored as immutable events |
| Soft Deletes | Regulatory data never physically deleted; `deleted_at` timestamp used |
| Audit Columns | Every table has `created_at`, `updated_at`, `created_by`, `updated_by` |
| UUID Primary Keys | All primary keys are UUIDv7 (time-ordered) for distributed generation and index performance |
| Timezone | All timestamps stored as `TIMESTAMPTZ` in UTC |

---

## 2. Entity Relationship Diagram

```
+==================+       +==================+       +==================+
|   AUTH_DB        |       |   ACCOUNT_DB     |       |   LEDGER_DB      |
|==================|       |==================|       |==================|
| users_credentials|       | users            |       | chart_of_accounts|
| sessions         |  ref  | accounts         |  ref  | journal_entries  |
| devices          |------>| sub_accounts     |------>| ledger_entries   |
| mfa_secrets      |       | beneficiaries    |       | event_store      |
| refresh_tokens   |       | account_tiers    |       | balances         |
+==================+       | user_preferences |       | reconciliations  |
                           +==================+       | fx_rates_history |
                                                      +==================+
                                    |
                    +---------------+---------------+
                    |                               |
          +=========v======+              +========v=======+
          |   PAYMENT_DB   |              |   CARD_DB      |
          |================|              |================|
          | payment_orders |              | cards          |
          | payment_routes |              | card_controls  |
          | scheduled_pay  |              | authorizations |
          | sdd_mandates   |              | disputes       |
          | fx_orders      |              | tokenizations  |
          +================+              +================+

          +================+              +================+
          |   CRYPTO_DB    |              |   KYC_DB       |
          |================|              |================|
          | wallets        |              | verifications  |
          | crypto_orders  |              | aml_screenings |
          | blockchain_txs |              | risk_scores    |
          | price_feeds    |              | review_queue   |
          +================+              | documents_meta |
                                          +================+

          +================+              +================+
          | NOTIFICATION_DB|              |   AUDIT_DB     |
          |================|              |================|
          | templates      |              | audit_events   |
          | delivery_log   |              | compliance_rpts|
          | preferences    |              | data_requests  |
          | device_tokens  |              | retention_rules|
          +================+              +================+

          +================+
          |   FRAUD_DB     |
          |================|
          | fraud_rules    |
          | fraud_cases    |
          | velocity_logs  |
          | risk_signals   |
          +================+
```

---

## 3. Schema Definitions

### 3.1 Auth Database (`auth_db`)

```sql
-- User credentials (separated from profile for security isolation)
CREATE TABLE user_credentials (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE,          -- References account_db.users.id
    email           VARCHAR(255) NOT NULL UNIQUE,
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    phone           VARCHAR(20) NOT NULL UNIQUE,
    phone_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    password_hash   TEXT NOT NULL,                  -- bcrypt/argon2id hash
    password_salt   TEXT NOT NULL,
    failed_attempts INT NOT NULL DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    last_login_at   TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'locked', 'suspended', 'closed')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Active sessions
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    device_id       UUID NOT NULL,
    access_token_jti VARCHAR(64) NOT NULL UNIQUE,   -- JWT ID for revocation check
    ip_address      INET,
    user_agent      TEXT,
    location        VARCHAR(100),
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Registered devices
CREATE TABLE devices (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    device_name     VARCHAR(100),
    device_type     VARCHAR(20) CHECK (device_type IN ('ios', 'android')),
    device_fingerprint TEXT NOT NULL,
    push_token      TEXT,                           -- APNs or FCM token
    biometric_key   TEXT,                           -- Public key for biometric auth
    is_trusted      BOOLEAN NOT NULL DEFAULT FALSE,
    registered_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at    TIMESTAMPTZ
);

CREATE INDEX idx_devices_user_id ON devices(user_id);

-- MFA secrets
CREATE TABLE mfa_secrets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    type            VARCHAR(20) NOT NULL CHECK (type IN ('totp', 'sms', 'push')),
    secret          TEXT,                           -- Encrypted TOTP secret
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    device_id       UUID NOT NULL REFERENCES devices(id),
    token_hash      VARCHAR(64) NOT NULL UNIQUE,    -- SHA-256 of actual token
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id, device_id);
```

### 3.2 Account Database (`account_db`)

```sql
-- Currency reference
CREATE TABLE currencies (
    code            VARCHAR(3) PRIMARY KEY,         -- ISO 4217
    name            VARCHAR(50) NOT NULL,
    decimal_places  INT NOT NULL DEFAULT 2,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

INSERT INTO currencies (code, name, decimal_places) VALUES
    ('EUR', 'Euro', 2),
    ('USD', 'US Dollar', 2),
    ('GBP', 'British Pound', 2),
    ('PLN', 'Polish Zloty', 2),
    ('CHF', 'Swiss Franc', 2);

-- Account tiers
CREATE TABLE account_tiers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(20) NOT NULL UNIQUE
                    CHECK (name IN ('basic', 'standard', 'premium')),
    monthly_fee     NUMERIC(19,4) NOT NULL DEFAULT 0,
    daily_transfer_limit    NUMERIC(19,4) NOT NULL,
    monthly_transfer_limit  NUMERIC(19,4) NOT NULL,
    daily_card_limit        NUMERIC(19,4) NOT NULL,
    monthly_atm_limit       NUMERIC(19,4) NOT NULL,
    free_atm_withdrawals    INT NOT NULL DEFAULT 0,
    fx_markup_percent       NUMERIC(5,4) NOT NULL,  -- e.g., 0.0050 = 0.50%
    max_sub_accounts        INT NOT NULL DEFAULT 5,
    features        JSONB NOT NULL DEFAULT '{}',    -- Feature flags
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO account_tiers (id, name, daily_transfer_limit, monthly_transfer_limit,
    daily_card_limit, monthly_atm_limit, free_atm_withdrawals, fx_markup_percent) VALUES
    (gen_random_uuid(), 'basic', 2000, 10000, 1000, 500, 2, 0.0050),
    (gen_random_uuid(), 'standard', 10000, 50000, 5000, 1500, 5, 0.0030),
    (gen_random_uuid(), 'premium', 50000, 200000, 25000, 5000, 10, 0.0010);

-- User profiles
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id     VARCHAR(20) NOT NULL UNIQUE,    -- Human-readable ID (TP-XXXXXXXX)
    tier_id         UUID NOT NULL REFERENCES account_tiers(id),
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    date_of_birth   DATE NOT NULL,
    nationality     VARCHAR(3),                     -- ISO 3166-1 alpha-3
    tax_residency   VARCHAR(3),
    tax_id          VARCHAR(50),                    -- Encrypted
    address_line1   VARCHAR(200),
    address_line2   VARCHAR(200),
    city            VARCHAR(100),
    postal_code     VARCHAR(20),
    country         VARCHAR(3) NOT NULL,            -- ISO 3166-1 alpha-3
    language        VARCHAR(5) NOT NULL DEFAULT 'en',
    kyc_status      VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (kyc_status IN ('pending', 'in_progress', 'verified',
                        'rejected', 'expired', 'recheck_required')),
    kyc_level       INT NOT NULL DEFAULT 0,         -- 0=none, 1=basic, 2=enhanced, 3=full
    risk_score      INT DEFAULT 0,                  -- 0-100
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'active', 'suspended', 'closed')),
    pep_status      BOOLEAN NOT NULL DEFAULT FALSE,
    sanctions_match BOOLEAN NOT NULL DEFAULT FALSE,
    migrated_from   VARCHAR(50),                    -- Legacy system reference
    closed_at       TIMESTAMPTZ,
    closed_reason   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ                     -- Soft delete
);

CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_kyc_status ON users(kyc_status);
CREATE INDEX idx_users_external_id ON users(external_id);

-- Accounts (one per user, container for sub-accounts)
CREATE TABLE accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    account_number  VARCHAR(34) NOT NULL UNIQUE,    -- Internal account number
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'frozen', 'closed')),
    opened_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);

-- Sub-accounts (one per currency)
CREATE TABLE sub_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL REFERENCES accounts(id),
    currency        VARCHAR(3) NOT NULL REFERENCES currencies(code),
    iban            VARCHAR(34) UNIQUE,             -- Lithuanian IBAN (LTxx...)
    bic             VARCHAR(11),                    -- SWIFT/BIC
    ledger_account_id UUID NOT NULL,                -- Reference to ledger_db.chart_of_accounts
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'frozen', 'closed')),
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(account_id, currency)
);

CREATE INDEX idx_sub_accounts_iban ON sub_accounts(iban);
CREATE INDEX idx_sub_accounts_ledger ON sub_accounts(ledger_account_id);

-- Saved beneficiaries / payees
CREATE TABLE beneficiaries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    name            VARCHAR(200) NOT NULL,
    iban            VARCHAR(34) NOT NULL,
    bic             VARCHAR(11),
    bank_name       VARCHAR(200),
    default_reference VARCHAR(140),
    is_internal     BOOLEAN NOT NULL DEFAULT FALSE,  -- TeslaPay-to-TeslaPay
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_beneficiaries_user ON beneficiaries(user_id) WHERE deleted_at IS NULL;

-- User notification preferences
CREATE TABLE user_preferences (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES users(id),
    push_transactions   BOOLEAN NOT NULL DEFAULT TRUE,
    push_security       BOOLEAN NOT NULL DEFAULT TRUE,
    push_marketing      BOOLEAN NOT NULL DEFAULT FALSE,
    email_transactions  BOOLEAN NOT NULL DEFAULT TRUE,
    email_security      BOOLEAN NOT NULL DEFAULT TRUE,
    email_marketing     BOOLEAN NOT NULL DEFAULT FALSE,
    sms_security        BOOLEAN NOT NULL DEFAULT TRUE,
    dark_mode           VARCHAR(10) NOT NULL DEFAULT 'system'
                        CHECK (dark_mode IN ('light', 'dark', 'system')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 3.3 Ledger Database (`ledger_db`) -- Core Double-Entry Bookkeeping

```sql
-- Chart of accounts
-- Follows standard accounting classification
CREATE TABLE chart_of_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(20) NOT NULL UNIQUE,    -- e.g., '1100-EUR-001234'
    name            VARCHAR(200) NOT NULL,
    type            VARCHAR(20) NOT NULL
                    CHECK (type IN ('asset', 'liability', 'equity', 'revenue', 'expense')),
    category        VARCHAR(30) NOT NULL
                    CHECK (category IN (
                        'customer_funds',           -- Liability: customer balances
                        'safeguarded_funds',        -- Asset: funds held at safeguarding bank
                        'fee_revenue',              -- Revenue: transaction/card/fx fees
                        'fx_revenue',               -- Revenue: FX spread income
                        'interest_expense',         -- Expense: interest paid to customers
                        'operational',              -- Expense: operational costs
                        'settlement',               -- Asset/Liability: pending settlements
                        'suspense',                 -- Temporary holding
                        'control'                   -- Control accounts
                    )),
    currency        VARCHAR(3) NOT NULL,
    parent_id       UUID REFERENCES chart_of_accounts(id),
    is_system       BOOLEAN NOT NULL DEFAULT FALSE, -- System accounts cannot be modified
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_coa_type ON chart_of_accounts(type);
CREATE INDEX idx_coa_category ON chart_of_accounts(category);
CREATE INDEX idx_coa_currency ON chart_of_accounts(currency);

-- Journal entries (the "header" for a set of balanced entries)
CREATE TABLE journal_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    posting_id      VARCHAR(64) NOT NULL UNIQUE,    -- Idempotency key
    effective_date  DATE NOT NULL,
    description     VARCHAR(500) NOT NULL,
    entry_type      VARCHAR(30) NOT NULL
                    CHECK (entry_type IN (
                        'payment_debit', 'payment_credit',
                        'internal_transfer',
                        'card_authorization', 'card_settlement', 'card_reversal',
                        'fx_exchange',
                        'crypto_buy', 'crypto_sell',
                        'fee', 'interest',
                        'adjustment', 'reversal',
                        'opening_balance', 'closing'
                    )),
    status          VARCHAR(20) NOT NULL DEFAULT 'posted'
                    CHECK (status IN ('pending', 'posted', 'reversed')),
    reference_type  VARCHAR(30),                    -- 'payment_order', 'card_auth', etc.
    reference_id    UUID,                           -- FK to source entity
    reversal_of     UUID REFERENCES journal_entries(id),
    metadata        JSONB DEFAULT '{}',
    created_by      VARCHAR(100) NOT NULL,          -- Service or user that created this
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_journal_effective_date ON journal_entries(effective_date);
CREATE INDEX idx_journal_type ON journal_entries(entry_type);
CREATE INDEX idx_journal_reference ON journal_entries(reference_type, reference_id);
CREATE INDEX idx_journal_posting_id ON journal_entries(posting_id);

-- Ledger entries (individual debit/credit lines within a journal entry)
-- This is the core double-entry table
CREATE TABLE ledger_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id),
    account_id      UUID NOT NULL REFERENCES chart_of_accounts(id),
    entry_side      VARCHAR(6) NOT NULL CHECK (entry_side IN ('debit', 'credit')),
    amount          NUMERIC(19,4) NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL,
    balance_after   NUMERIC(19,4) NOT NULL,         -- Running balance after this entry
    sequence_num    BIGINT NOT NULL,                 -- Per-account sequence for ordering
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Monthly partitions (created dynamically by migration scripts)
CREATE TABLE ledger_entries_2026_03 PARTITION OF ledger_entries
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE ledger_entries_2026_04 PARTITION OF ledger_entries
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
-- ... additional partitions created monthly by automation

CREATE INDEX idx_ledger_journal ON ledger_entries(journal_entry_id);
CREATE INDEX idx_ledger_account ON ledger_entries(account_id, sequence_num DESC);
CREATE INDEX idx_ledger_account_date ON ledger_entries(account_id, created_at DESC);

-- CRITICAL CONSTRAINT: Every journal entry must balance (debits = credits)
-- Enforced at application level via transaction:
--   BEGIN;
--   INSERT INTO journal_entries ...;
--   INSERT INTO ledger_entries (debit side) ...;
--   INSERT INTO ledger_entries (credit side) ...;
--   -- Verify: SELECT SUM(CASE WHEN entry_side='debit' THEN amount ELSE -amount END)
--   --         FROM ledger_entries WHERE journal_entry_id = $1;
--   -- Must equal 0, otherwise ROLLBACK
--   COMMIT;

-- Materialized balance view (CQRS read model)
CREATE TABLE account_balances (
    account_id      UUID PRIMARY KEY REFERENCES chart_of_accounts(id),
    currency        VARCHAR(3) NOT NULL,
    available       NUMERIC(19,4) NOT NULL DEFAULT 0,   -- Available for spending
    pending         NUMERIC(19,4) NOT NULL DEFAULT 0,   -- Authorized but not settled
    reserved        NUMERIC(19,4) NOT NULL DEFAULT 0,   -- Regulatory holds
    total           NUMERIC(19,4) NOT NULL DEFAULT 0,   -- available + pending + reserved
    last_entry_id   UUID,                                -- Last processed ledger entry
    last_sequence   BIGINT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version         BIGINT NOT NULL DEFAULT 0            -- Optimistic locking
);

-- Event store (append-only, for event sourcing)
CREATE TABLE event_store (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id    UUID NOT NULL,
    aggregate_type  VARCHAR(30) NOT NULL,
    event_type      VARCHAR(50) NOT NULL,
    event_data      JSONB NOT NULL,
    metadata        JSONB DEFAULT '{}',
    sequence_number BIGINT NOT NULL,
    checksum        VARCHAR(64) NOT NULL,            -- SHA-256 chain
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(aggregate_id, sequence_number)
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_event_store_aggregate ON event_store(aggregate_id, sequence_number);
CREATE INDEX idx_event_store_type ON event_store(event_type, created_at);

-- Reconciliation records
CREATE TABLE reconciliations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reconciliation_date DATE NOT NULL,
    type            VARCHAR(30) NOT NULL
                    CHECK (type IN ('daily_balance', 'sepa_settlement',
                        'card_settlement', 'safeguarding')),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'matched', 'discrepancy', 'resolved')),
    expected_amount NUMERIC(19,4) NOT NULL,
    actual_amount   NUMERIC(19,4) NOT NULL,
    difference      NUMERIC(19,4) NOT NULL DEFAULT 0,
    currency        VARCHAR(3) NOT NULL,
    notes           TEXT,
    resolved_by     VARCHAR(100),
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- FX rates history
CREATE TABLE fx_rates_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_currency VARCHAR(3) NOT NULL,
    target_currency VARCHAR(3) NOT NULL,
    mid_rate        NUMERIC(19,8) NOT NULL,          -- Mid-market rate
    buy_rate        NUMERIC(19,8) NOT NULL,          -- TeslaPay buy rate (with markup)
    sell_rate       NUMERIC(19,8) NOT NULL,          -- TeslaPay sell rate (with markup)
    provider        VARCHAR(30) NOT NULL,             -- ECB, aggregator, etc.
    valid_from      TIMESTAMPTZ NOT NULL,
    valid_to        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fx_rates_pair ON fx_rates_history(source_currency, target_currency, valid_from DESC);
```

### 3.4 Payment Database (`payment_db`)

```sql
-- Payment orders
CREATE TABLE payment_orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(64) NOT NULL UNIQUE,
    user_id         UUID NOT NULL,
    type            VARCHAR(20) NOT NULL
                    CHECK (type IN ('sepa_sct', 'sepa_sct_inst', 'sepa_sdd',
                        'internal', 'fx_exchange')),
    status          VARCHAR(20) NOT NULL DEFAULT 'created'
                    CHECK (status IN ('created', 'validating', 'authorized',
                        'submitted', 'processing', 'settled', 'completed',
                        'failed', 'cancelled', 'returned')),
    -- Source
    source_sub_account_id   UUID NOT NULL,
    source_iban             VARCHAR(34),
    source_currency         VARCHAR(3) NOT NULL,
    source_amount           NUMERIC(19,4) NOT NULL,
    -- Destination
    dest_sub_account_id     UUID,                    -- For internal transfers
    dest_iban               VARCHAR(34),
    dest_name               VARCHAR(200),
    dest_bic                VARCHAR(11),
    dest_currency           VARCHAR(3) NOT NULL,
    dest_amount             NUMERIC(19,4),            -- Different from source if FX
    -- FX
    fx_rate                 NUMERIC(19,8),
    fx_markup               NUMERIC(19,4),
    -- Payment details
    reference               VARCHAR(140),
    end_to_end_id           VARCHAR(35),             -- SEPA end-to-end reference
    -- Fees
    fee_amount              NUMERIC(19,4) NOT NULL DEFAULT 0,
    fee_currency            VARCHAR(3),
    -- External processing
    external_id             VARCHAR(100),             -- Banking Circle reference
    external_status         VARCHAR(50),
    -- Scheduling
    scheduled_date          DATE,
    recurring_id            UUID,                     -- Link to recurring schedule
    -- Tracking
    journal_entry_id        UUID,                     -- Link to ledger
    submitted_at            TIMESTAMPTZ,
    settled_at              TIMESTAMPTZ,
    failed_reason           TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_payment_orders_user ON payment_orders(user_id, created_at DESC);
CREATE INDEX idx_payment_orders_status ON payment_orders(status);
CREATE INDEX idx_payment_orders_external ON payment_orders(external_id);
CREATE INDEX idx_payment_orders_idem ON payment_orders(idempotency_key);

-- Scheduled / recurring payments
CREATE TABLE scheduled_payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    frequency       VARCHAR(20) NOT NULL
                    CHECK (frequency IN ('daily', 'weekly', 'biweekly',
                        'monthly', 'quarterly', 'yearly', 'custom')),
    custom_interval_days INT,                        -- For custom frequency
    next_execution  DATE NOT NULL,
    end_date        DATE,
    dest_iban       VARCHAR(34) NOT NULL,
    dest_name       VARCHAR(200) NOT NULL,
    amount          NUMERIC(19,4) NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    reference       VARCHAR(140),
    payment_type    VARCHAR(20) NOT NULL DEFAULT 'sepa_sct',
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'paused', 'completed', 'cancelled')),
    last_executed   DATE,
    failure_count   INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_scheduled_next ON scheduled_payments(next_execution)
    WHERE status = 'active';

-- SEPA Direct Debit mandates
CREATE TABLE sdd_mandates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    mandate_id      VARCHAR(35) NOT NULL UNIQUE,     -- SEPA mandate reference
    creditor_name   VARCHAR(200) NOT NULL,
    creditor_iban   VARCHAR(34) NOT NULL,
    creditor_id     VARCHAR(35) NOT NULL,            -- Creditor identifier
    scheme          VARCHAR(10) NOT NULL
                    CHECK (scheme IN ('CORE', 'B2B')),
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'suspended', 'cancelled')),
    signed_date     DATE NOT NULL,
    last_collection DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- FX exchange orders
CREATE TABLE fx_orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    source_currency VARCHAR(3) NOT NULL,
    target_currency VARCHAR(3) NOT NULL,
    source_amount   NUMERIC(19,4) NOT NULL,
    target_amount   NUMERIC(19,4) NOT NULL,
    rate            NUMERIC(19,8) NOT NULL,
    mid_market_rate NUMERIC(19,8) NOT NULL,
    markup          NUMERIC(5,4) NOT NULL,
    rate_locked_at  TIMESTAMPTZ NOT NULL,
    rate_expires_at TIMESTAMPTZ NOT NULL,            -- 30-second lock
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'executed', 'expired', 'failed')),
    journal_entry_id UUID,
    executed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 3.5 Card Database (`card_db`)

```sql
-- Cards (no PAN/CVV stored -- those stay at Enfuce)
CREATE TABLE cards (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    account_id      UUID NOT NULL,
    sub_account_id  UUID NOT NULL,                   -- Linked currency sub-account
    processor_card_id VARCHAR(100) NOT NULL UNIQUE,  -- Enfuce card reference
    card_type       VARCHAR(10) NOT NULL
                    CHECK (card_type IN ('virtual', 'physical')),
    card_brand      VARCHAR(20) NOT NULL DEFAULT 'mastercard',
    last_four       VARCHAR(4) NOT NULL,
    expiry_month    INT NOT NULL,
    expiry_year     INT NOT NULL,
    cardholder_name VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'inactive'
                    CHECK (status IN ('inactive', 'active', 'frozen',
                        'blocked', 'expired', 'cancelled')),
    -- Physical card specifics
    delivery_address JSONB,
    tracking_number VARCHAR(50),
    delivered_at    TIMESTAMPTZ,
    -- Replacement
    replaces_card_id UUID REFERENCES cards(id),
    replacement_reason VARCHAR(30),
    -- Timestamps
    activated_at    TIMESTAMPTZ,
    frozen_at       TIMESTAMPTZ,
    blocked_at      TIMESTAMPTZ,
    blocked_reason  VARCHAR(30),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cards_user ON cards(user_id);
CREATE INDEX idx_cards_processor ON cards(processor_card_id);

-- Card spending controls
CREATE TABLE card_controls (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id         UUID NOT NULL REFERENCES cards(id),
    per_transaction_limit   NUMERIC(19,4),
    daily_limit             NUMERIC(19,4),
    monthly_limit           NUMERIC(19,4),
    atm_daily_limit         NUMERIC(19,4),
    online_enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    contactless_enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    atm_enabled             BOOLEAN NOT NULL DEFAULT TRUE,
    magstripe_enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    blocked_mcc_codes       TEXT[],                  -- Merchant Category Codes
    allowed_countries       TEXT[],                  -- ISO 3166 country codes (empty = all)
    blocked_countries       TEXT[],
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Card authorizations
CREATE TABLE authorizations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id         UUID NOT NULL REFERENCES cards(id),
    processor_auth_id VARCHAR(100) NOT NULL,
    type            VARCHAR(20) NOT NULL
                    CHECK (type IN ('purchase', 'atm', 'contactless',
                        'online', 'recurring', 'refund')),
    status          VARCHAR(20) NOT NULL
                    CHECK (status IN ('authorized', 'declined', 'settled',
                        'reversed', 'expired')),
    amount          NUMERIC(19,4) NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    billing_amount  NUMERIC(19,4),                   -- In card currency
    billing_currency VARCHAR(3),
    merchant_name   VARCHAR(200),
    merchant_mcc    VARCHAR(4),
    merchant_country VARCHAR(3),
    merchant_city   VARCHAR(100),
    pos_entry_mode  VARCHAR(20),
    decline_reason  VARCHAR(50),
    three_ds_result VARCHAR(20),
    is_apple_pay    BOOLEAN NOT NULL DEFAULT FALSE,
    is_google_pay   BOOLEAN NOT NULL DEFAULT FALSE,
    journal_entry_id UUID,
    authorized_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settled_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_auth_card ON authorizations(card_id, created_at DESC);
CREATE INDEX idx_auth_processor ON authorizations(processor_auth_id);

-- Card disputes
CREATE TABLE disputes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id         UUID NOT NULL REFERENCES cards(id),
    authorization_id UUID NOT NULL,
    reason          VARCHAR(30) NOT NULL
                    CHECK (reason IN ('unauthorized', 'duplicate',
                        'goods_not_received', 'amount_incorrect', 'other')),
    description     TEXT,
    evidence_urls   TEXT[],
    amount          NUMERIC(19,4) NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'submitted'
                    CHECK (status IN ('submitted', 'investigating',
                        'resolved_for_user', 'resolved_for_merchant', 'closed')),
    processor_dispute_id VARCHAR(100),
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tokenizations (Apple Pay / Google Pay)
CREATE TABLE tokenizations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id         UUID NOT NULL REFERENCES cards(id),
    wallet_type     VARCHAR(20) NOT NULL
                    CHECK (wallet_type IN ('apple_pay', 'google_pay', 'samsung_pay')),
    token_reference VARCHAR(100) NOT NULL,           -- MDES token reference
    device_name     VARCHAR(100),
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'suspended', 'deleted')),
    provisioned_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_tokenizations_card ON tokenizations(card_id);
```

### 3.6 Crypto Database (`crypto_db`)

```sql
-- Fuse Smart Wallets
CREATE TABLE wallets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE,
    smart_wallet_address VARCHAR(42) NOT NULL UNIQUE, -- Fuse Smart Wallet address (0x...)
    eoa_address     VARCHAR(42) NOT NULL UNIQUE,      -- Externally Owned Account (owner)
    status          VARCHAR(20) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('creating', 'active', 'frozen', 'closed')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- NOTE: Private keys are NEVER stored server-side. The EOA private key is
-- derived from device-local secure storage (Keychain/Keystore). The Smart
-- Wallet is a contract wallet controlled by the EOA.

-- Crypto buy/sell orders (on-ramp / off-ramp)
CREATE TABLE crypto_orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(64) NOT NULL UNIQUE,
    user_id         UUID NOT NULL,
    wallet_id       UUID NOT NULL REFERENCES wallets(id),
    type            VARCHAR(10) NOT NULL CHECK (type IN ('buy', 'sell')),
    -- Fiat side
    fiat_amount     NUMERIC(19,4) NOT NULL,
    fiat_currency   VARCHAR(3) NOT NULL DEFAULT 'EUR',
    -- Crypto side
    token_symbol    VARCHAR(10) NOT NULL,            -- FUSE, USDC, USDT
    token_amount    NUMERIC(28,18) NOT NULL,         -- 18 decimal places for ERC-20
    token_address   VARCHAR(42) NOT NULL,            -- Token contract address
    -- Pricing
    exchange_rate   NUMERIC(28,18) NOT NULL,
    fee_amount      NUMERIC(19,4) NOT NULL,
    fee_currency    VARCHAR(3) NOT NULL,
    -- Execution
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'executing', 'completed',
                        'failed', 'cancelled')),
    blockchain_tx_hash VARCHAR(66),                  -- 0x...
    journal_entry_id UUID,                           -- Fiat leg in ledger
    error_message   TEXT,
    executed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_crypto_orders_user ON crypto_orders(user_id, created_at DESC);

-- Blockchain transactions (sends/receives)
CREATE TABLE blockchain_transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id       UUID NOT NULL REFERENCES wallets(id),
    direction       VARCHAR(10) NOT NULL CHECK (direction IN ('incoming', 'outgoing')),
    tx_hash         VARCHAR(66) NOT NULL UNIQUE,
    block_number    BIGINT,
    from_address    VARCHAR(42) NOT NULL,
    to_address      VARCHAR(42) NOT NULL,
    token_symbol    VARCHAR(10) NOT NULL,
    token_address   VARCHAR(42),
    amount          NUMERIC(28,18) NOT NULL,
    gas_used        BIGINT,
    gas_price       NUMERIC(28,18),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'confirmed', 'failed')),
    confirmations   INT NOT NULL DEFAULT 0,
    fiat_value_eur  NUMERIC(19,4),                   -- EUR value at time of tx
    detected_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blockchain_tx_wallet ON blockchain_transactions(wallet_id, created_at DESC);
CREATE INDEX idx_blockchain_tx_hash ON blockchain_transactions(tx_hash);

-- Price feeds (cached from aggregator)
CREATE TABLE price_feeds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_symbol    VARCHAR(10) NOT NULL,
    price_eur       NUMERIC(19,8) NOT NULL,
    price_usd       NUMERIC(19,8) NOT NULL,
    change_24h_pct  NUMERIC(8,4),
    volume_24h_eur  NUMERIC(19,4),
    source          VARCHAR(30) NOT NULL,
    fetched_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_price_feeds_token ON price_feeds(token_symbol, fetched_at DESC);
```

### 3.7 KYC Database (`kyc_db`)

```sql
-- Verification records
CREATE TABLE verifications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    sumsub_applicant_id VARCHAR(100) NOT NULL UNIQUE,
    level           VARCHAR(30) NOT NULL
                    CHECK (level IN ('basic', 'enhanced', 'full')),
    status          VARCHAR(20) NOT NULL DEFAULT 'initiated'
                    CHECK (status IN ('initiated', 'pending', 'approved',
                        'rejected', 'retry', 'expired')),
    -- Document details (metadata only -- docs stored in Sumsub)
    document_type   VARCHAR(30),
    document_country VARCHAR(3),
    document_number_masked VARCHAR(20),              -- Last 4 characters only
    -- Verification results
    liveness_passed BOOLEAN,
    nfc_verified    BOOLEAN DEFAULT FALSE,
    review_result   VARCHAR(10),                     -- GREEN, RED
    reject_type     VARCHAR(10),                     -- FINAL, RETRY
    reject_labels   TEXT[],
    -- Risk
    risk_score      INT,
    risk_level      VARCHAR(10) CHECK (risk_level IN ('low', 'medium', 'high')),
    -- Tracking
    initiated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    attempt_number  INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_verifications_user ON verifications(user_id);
CREATE INDEX idx_verifications_sumsub ON verifications(sumsub_applicant_id);
CREATE INDEX idx_verifications_status ON verifications(status);

-- AML screening results
CREATE TABLE aml_screenings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    verification_id UUID REFERENCES verifications(id),
    type            VARCHAR(20) NOT NULL
                    CHECK (type IN ('onboarding', 'ongoing', 'triggered')),
    -- Results
    sanctions_match BOOLEAN NOT NULL DEFAULT FALSE,
    pep_match       BOOLEAN NOT NULL DEFAULT FALSE,
    adverse_media   BOOLEAN NOT NULL DEFAULT FALSE,
    match_details   JSONB,                           -- Detailed match information
    -- Review
    review_status   VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (review_status IN ('pending', 'auto_cleared',
                        'manual_review', 'cleared', 'confirmed_match', 'escalated')),
    reviewed_by     VARCHAR(100),
    review_notes    TEXT,
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aml_user ON aml_screenings(user_id, created_at DESC);

-- Risk scoring history
CREATE TABLE risk_scores (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    score           INT NOT NULL CHECK (score BETWEEN 0 AND 100),
    factors         JSONB NOT NULL,                  -- Contributing factors
    previous_score  INT,
    trigger_event   VARCHAR(50) NOT NULL,            -- What caused re-scoring
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Manual review queue
CREATE TABLE review_queue (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    type            VARCHAR(30) NOT NULL
                    CHECK (type IN ('kyc_review', 'aml_alert', 'source_of_funds',
                        'proof_of_address', 'recheck', 'sar_review')),
    priority        INT NOT NULL DEFAULT 5,          -- 1=highest, 10=lowest
    status          VARCHAR(20) NOT NULL DEFAULT 'open'
                    CHECK (status IN ('open', 'assigned', 'in_review',
                        'approved', 'rejected', 'escalated')),
    assigned_to     VARCHAR(100),
    data            JSONB NOT NULL,                  -- Relevant data for review
    decision        VARCHAR(20),
    decision_reason TEXT,
    decided_at      TIMESTAMPTZ,
    sla_deadline    TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_queue_status ON review_queue(status, priority, created_at);
```

### 3.8 Audit Database (`audit_db`)

```sql
-- Immutable audit event log
CREATE TABLE audit_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type      VARCHAR(50) NOT NULL,
    actor_type      VARCHAR(20) NOT NULL
                    CHECK (actor_type IN ('user', 'system', 'admin', 'external')),
    actor_id        VARCHAR(100) NOT NULL,
    target_type     VARCHAR(30) NOT NULL,            -- 'user', 'account', 'payment', etc.
    target_id       UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,            -- 'created', 'updated', 'deleted', etc.
    changes         JSONB,                           -- {field: {old: x, new: y}}
    metadata        JSONB DEFAULT '{}',              -- correlation_id, ip, device, etc.
    ip_address      INET,
    user_agent      TEXT,
    checksum        VARCHAR(64) NOT NULL,            -- SHA-256 chain for tamper detection
    previous_checksum VARCHAR(64),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Quarterly partitions (7-year retention = 28 partitions)
CREATE TABLE audit_events_2026_q1 PARTITION OF audit_events
    FOR VALUES FROM ('2026-01-01') TO ('2026-04-01');
CREATE TABLE audit_events_2026_q2 PARTITION OF audit_events
    FOR VALUES FROM ('2026-04-01') TO ('2026-07-01');

CREATE INDEX idx_audit_target ON audit_events(target_type, target_id, created_at DESC);
CREATE INDEX idx_audit_actor ON audit_events(actor_id, created_at DESC);
CREATE INDEX idx_audit_type ON audit_events(event_type, created_at DESC);

-- NOTE: This table is append-only. No UPDATE or DELETE operations are permitted.
-- Database user for Audit Service has only INSERT and SELECT privileges.
REVOKE UPDATE, DELETE ON audit_events FROM audit_service_user;

-- GDPR data subject requests
CREATE TABLE data_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    type            VARCHAR(20) NOT NULL
                    CHECK (type IN ('access', 'correction', 'deletion', 'portability')),
    status          VARCHAR(20) NOT NULL DEFAULT 'received'
                    CHECK (status IN ('received', 'processing', 'completed',
                        'partially_completed', 'denied')),
    description     TEXT,
    response_data   JSONB,                           -- For access/portability requests
    exceptions      TEXT[],                          -- Retention exceptions for deletion
    completed_at    TIMESTAMPTZ,
    deadline        TIMESTAMPTZ NOT NULL,            -- GDPR 30-day deadline
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Data retention rules
CREATE TABLE retention_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_category   VARCHAR(50) NOT NULL UNIQUE,
    retention_years INT NOT NULL,
    legal_basis     TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO retention_rules (id, data_category, retention_years, legal_basis) VALUES
    (gen_random_uuid(), 'transaction_data', 5, 'EMD2/PSD2 regulatory requirement'),
    (gen_random_uuid(), 'kyc_documents', 5, 'AMLD6 Art. 40'),
    (gen_random_uuid(), 'audit_trail', 7, 'Internal audit policy + regulatory'),
    (gen_random_uuid(), 'customer_profile', 5, 'AMLD6 Art. 40, post account closure'),
    (gen_random_uuid(), 'card_transaction_data', 5, 'PSD2 regulatory requirement'),
    (gen_random_uuid(), 'communication_records', 5, 'PSD2 regulatory requirement');
```

### 3.9 Notification Database (`notification_db`)

```sql
-- Notification templates
CREATE TABLE templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL UNIQUE,
    channel         VARCHAR(10) NOT NULL
                    CHECK (channel IN ('push', 'email', 'sms', 'in_app')),
    language        VARCHAR(5) NOT NULL DEFAULT 'en',
    subject         VARCHAR(200),                    -- For email
    body_template   TEXT NOT NULL,                   -- Handlebars template
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name, channel, language)
);

-- Delivery log
CREATE TABLE delivery_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    template_id     UUID REFERENCES templates(id),
    channel         VARCHAR(10) NOT NULL,
    recipient       VARCHAR(200) NOT NULL,           -- Email, phone, or device token
    subject         VARCHAR(200),
    body            TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'queued'
                    CHECK (status IN ('queued', 'sent', 'delivered',
                        'failed', 'bounced')),
    external_id     VARCHAR(100),                    -- Provider message ID
    error_message   TEXT,
    sent_at         TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_delivery_user ON delivery_log(user_id, created_at DESC);
```

### 3.10 Fraud Detection Database (`fraud_db`)

```sql
-- Fraud detection rules
CREATE TABLE fraud_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    rule_type       VARCHAR(20) NOT NULL
                    CHECK (rule_type IN ('velocity', 'amount', 'geographic',
                        'pattern', 'blacklist', 'whitelist')),
    conditions      JSONB NOT NULL,                  -- Rule conditions in JSON
    action          VARCHAR(20) NOT NULL
                    CHECK (action IN ('flag', 'decline', 'hold', 'alert')),
    severity        VARCHAR(10) NOT NULL
                    CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Fraud cases
CREATE TABLE fraud_cases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    triggered_rule_id UUID REFERENCES fraud_rules(id),
    transaction_type VARCHAR(20) NOT NULL,
    transaction_id  UUID NOT NULL,
    risk_score      INT NOT NULL CHECK (risk_score BETWEEN 0 AND 100),
    signals         JSONB NOT NULL,                  -- Risk signals that triggered the case
    status          VARCHAR(20) NOT NULL DEFAULT 'open'
                    CHECK (status IN ('open', 'investigating', 'confirmed_fraud',
                        'false_positive', 'closed')),
    assigned_to     VARCHAR(100),
    resolution_notes TEXT,
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fraud_cases_status ON fraud_cases(status, created_at DESC);
CREATE INDEX idx_fraud_cases_user ON fraud_cases(user_id);
```

---

## 4. Double-Entry Bookkeeping Model

### 4.1 Account Types and Normal Balances

```
Type         | Normal Balance | Debit Effect | Credit Effect
-------------|---------------|--------------|---------------
Asset        | Debit         | Increase     | Decrease
Liability    | Credit        | Decrease     | Increase
Equity       | Credit        | Decrease     | Increase
Revenue      | Credit        | Decrease     | Increase
Expense      | Debit         | Increase     | Decrease
```

### 4.2 Standard Posting Patterns

**Internal Transfer (User A -> User B, EUR 100):**
```
Journal Entry: "Internal transfer TP-A to TP-B"
  DEBIT   Customer Funds (User A EUR)       100.00
  CREDIT  Customer Funds (User B EUR)       100.00
```

**SEPA Outgoing Payment (EUR 500):**
```
Journal Entry: "SEPA SCT to DE89..."
  DEBIT   Customer Funds (User EUR)         500.00
  CREDIT  SEPA Settlement Pending           500.00

-- After settlement confirmation:
  DEBIT   SEPA Settlement Pending           500.00
  CREDIT  Safeguarded Funds (Bank Account)  500.00
```

**SEPA Incoming Payment (EUR 1,000):**
```
Journal Entry: "SEPA SCT from LT12..."
  DEBIT   Safeguarded Funds (Bank Account)  1000.00
  CREDIT  Customer Funds (User EUR)         1000.00
```

**Card Purchase Authorization (EUR 50):**
```
Journal Entry: "Card auth at Amazon.de"
  DEBIT   Customer Funds (User EUR)         50.00
  CREDIT  Card Authorization Holds          50.00

-- After settlement:
  DEBIT   Card Authorization Holds          50.00
  CREDIT  Safeguarded Funds (Bank Account)  50.00
```

**FX Exchange (EUR 100 -> USD at 1.0850):**
```
Journal Entry: "FX EUR/USD"
  DEBIT   Customer Funds (User EUR)         100.00
  CREDIT  Customer Funds (User USD)         108.50
  CREDIT  FX Revenue                        0.50    -- Markup
  DEBIT   FX Settlement                     108.00  -- At mid-market
```

**Crypto Buy (EUR 200 -> USDC):**
```
Journal Entry: "Crypto buy USDC"
  DEBIT   Customer Funds (User EUR)         200.00
  CREDIT  Crypto Settlement                 198.00  -- Net of fee
  CREDIT  Crypto Fee Revenue                2.00
```

### 4.3 Balance Invariant

The system enforces this invariant at all times:
```
SUM(all debits) = SUM(all credits)
```

This is verified:
1. **Per journal entry:** Application-level check before COMMIT
2. **Per account:** Running balance_after is validated against sum of entries
3. **Daily reconciliation:** Automated job verifies global balance equation
4. **Monitoring alert:** Any imbalance triggers a critical P1 alert

---

## 5. Partitioning Strategy

### 5.1 Time-Based Partitioning

| Table | Partition Granularity | Retention | Rationale |
|-------|----------------------|-----------|-----------|
| `ledger_entries` | Monthly | 5 years (60 partitions) | High write volume; older partitions rarely queried |
| `event_store` | Monthly | 7 years (84 partitions) | Event replay needs; long-term audit |
| `payment_orders` | Monthly | 5 years | Query patterns are date-bounded |
| `authorizations` | Monthly | 5 years | Card transaction volume |
| `audit_events` | Quarterly | 7 years (28 partitions) | Lower volume but longest retention |
| `delivery_log` | Monthly | 1 year (12 partitions) | High volume, low retention need |

### 5.2 Partition Management

```sql
-- Automated partition creation (run monthly by cron job)
-- Example for ledger_entries:
DO $$
DECLARE
    partition_date DATE := DATE_TRUNC('month', NOW() + INTERVAL '1 month');
    partition_name TEXT;
    start_date TEXT;
    end_date TEXT;
BEGIN
    partition_name := 'ledger_entries_' || TO_CHAR(partition_date, 'YYYY_MM');
    start_date := TO_CHAR(partition_date, 'YYYY-MM-DD');
    end_date := TO_CHAR(partition_date + INTERVAL '1 month', 'YYYY-MM-DD');

    EXECUTE FORMAT(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF ledger_entries
         FOR VALUES FROM (%L) TO (%L)',
        partition_name, start_date, end_date
    );
END $$;
```

### 5.3 Old Partition Archival

Partitions older than the active query window (e.g., 2 years for ledger_entries) are detached and moved to cold storage (S3 via `pg_dump`). They remain queryable by re-attaching if needed for regulatory requests.

---

## 6. Migration Strategy

### 6.1 Tool

All migrations managed by `golang-migrate` with versioned SQL files:

```
migrations/
  000001_create_users.up.sql
  000001_create_users.down.sql
  000002_create_accounts.up.sql
  000002_create_accounts.down.sql
  ...
```

### 6.2 Migration Rules

1. **Forward-only in production:** Down migrations exist for development rollback but are never used in production. Fixes are forward migrations.
2. **Backward compatible:** Every migration must be compatible with the currently running application version (enable blue-green deployments).
3. **No data loss:** Migrations that drop columns or tables must first verify the column/table is unused.
4. **Reviewed:** All migrations require approval from at least one database engineer.
5. **Tested:** Migrations tested against production-like dataset in staging before production deployment.

### 6.3 Legacy Data Migration

For migrating existing TeslaPay customers (PRD user story US-1.7):

1. **Extract** from legacy system via read-only database connection
2. **Transform** to new schema format in a staging area
3. **Validate** data integrity (balances, account status)
4. **Load** in batches (1,000 users per batch) with idempotency keys
5. **Reconcile** old and new balances; manual review for discrepancies
6. **Cutover** with maintenance window (estimated 2-4 hours for ~10K users)

---

## 7. Indexing Strategy

### 7.1 Primary Query Patterns and Supporting Indexes

| Query Pattern | Table | Index | Columns |
|--------------|-------|-------|---------|
| Get user by external ID | users | `idx_users_external_id` | `external_id` |
| Get sub-account by IBAN | sub_accounts | `idx_sub_accounts_iban` | `iban` |
| Get balance by account | account_balances | PRIMARY KEY | `account_id` |
| List ledger entries by account | ledger_entries | `idx_ledger_account_date` | `account_id, created_at DESC` |
| Get payment by idempotency key | payment_orders | `idx_payment_orders_idem` | `idempotency_key` |
| List transactions by user + date | payment_orders | `idx_payment_orders_user` | `user_id, created_at DESC` |
| Get card by processor ID | cards | `idx_cards_processor` | `processor_card_id` |
| Get auth by processor ID | authorizations | `idx_auth_processor` | `processor_auth_id` |
| List audit by target | audit_events | `idx_audit_target` | `target_type, target_id, created_at DESC` |

### 7.2 Index Maintenance

- `REINDEX CONCURRENTLY` scheduled weekly for high-write tables
- `pg_stat_user_indexes` monitored for unused indexes (remove after 30 days unused)
- Bloat monitored via `pgstattuple`; `VACUUM (FULL)` during maintenance windows if needed

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CTO | TBD | | Pending |
| Database Engineer | TBD | | Pending |
| Principal Architect | Dream Team Architect | 2026-03-03 | Submitted |
