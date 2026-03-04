CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_account_id UUID NOT NULL,
    recipient_account_id UUID,
    recipient_iban VARCHAR(34),
    recipient_name VARCHAR(255),
    amount DECIMAL(20,4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    reference VARCHAR(35),
    description TEXT,
    idempotency_key VARCHAR(64) UNIQUE,
    fee_amount DECIMAL(20,4) DEFAULT 0,
    fee_currency VARCHAR(3) DEFAULT 'EUR',
    fx_rate DECIMAL(20,8),
    fx_from_currency VARCHAR(3),
    fx_to_currency VARCHAR(3),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_sender ON payments(sender_account_id);
CREATE INDEX IF NOT EXISTS idx_payments_recipient ON payments(recipient_account_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_idempotency ON payments(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_payments_created ON payments(created_at DESC);

CREATE TABLE IF NOT EXISTS scheduled_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    recipient_account_id UUID,
    recipient_iban VARCHAR(34),
    recipient_name VARCHAR(255),
    amount DECIMAL(20,4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    type VARCHAR(20) NOT NULL,
    schedule_type VARCHAR(20) NOT NULL,
    reference VARCHAR(35),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    next_execution TIMESTAMPTZ,
    last_execution TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scheduled_payments_account ON scheduled_payments(account_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_payments_next ON scheduled_payments(next_execution) WHERE is_active = true;
