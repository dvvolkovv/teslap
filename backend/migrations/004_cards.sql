-- 004_cards.sql: Card issuance and lifecycle management.
-- Creates cards and card_transactions tables.

CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    sub_account_id UUID,
    card_number_encrypted TEXT NOT NULL,
    last_four VARCHAR(4) NOT NULL,
    expiry_month INTEGER NOT NULL,
    expiry_year INTEGER NOT NULL,
    cvv_hash VARCHAR(128) NOT NULL,
    cardholder_name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    daily_limit DECIMAL(20,4) DEFAULT 5000,
    monthly_limit DECIMAL(20,4) DEFAULT 25000,
    daily_spent DECIMAL(20,4) DEFAULT 0,
    monthly_spent DECIMAL(20,4) DEFAULT 0,
    is_contactless BOOLEAN DEFAULT true,
    is_online BOOLEAN DEFAULT true,
    is_atm BOOLEAN DEFAULT true,
    allowed_countries TEXT[],
    blocked_mcc_codes INTEGER[],
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cards_account ON cards(account_id);
CREATE INDEX IF NOT EXISTS idx_cards_status ON cards(status);

CREATE TABLE IF NOT EXISTS card_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES cards(id),
    merchant_name VARCHAR(255),
    merchant_mcc INTEGER,
    amount DECIMAL(20,4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    original_amount DECIMAL(20,4),
    original_currency VARCHAR(3),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    type VARCHAR(20) NOT NULL,
    country VARCHAR(3),
    authorization_code VARCHAR(20),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_card_tx_card ON card_transactions(card_id);
CREATE INDEX IF NOT EXISTS idx_card_tx_created ON card_transactions(created_at DESC);
