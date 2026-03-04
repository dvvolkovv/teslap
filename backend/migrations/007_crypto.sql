CREATE TABLE IF NOT EXISTS crypto_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    address VARCHAR(42) NOT NULL,
    network VARCHAR(20) NOT NULL DEFAULT 'fuse',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crypto_wallet_user ON crypto_wallets(user_id);

CREATE TABLE IF NOT EXISTS crypto_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES crypto_wallets(id),
    token_symbol VARCHAR(10) NOT NULL,
    token_name VARCHAR(50) NOT NULL,
    token_address VARCHAR(42),
    balance DECIMAL(30,18) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(wallet_id, token_symbol)
);
CREATE INDEX idx_crypto_balance_wallet ON crypto_balances(wallet_id);

CREATE TABLE IF NOT EXISTS crypto_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES crypto_wallets(id),
    tx_hash VARCHAR(66),
    type VARCHAR(20) NOT NULL,
    token_symbol VARCHAR(10) NOT NULL,
    amount DECIMAL(30,18) NOT NULL,
    fiat_amount DECIMAL(20,4),
    fiat_currency VARCHAR(3),
    rate DECIMAL(20,8),
    fee_amount DECIMAL(20,4),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    recipient_address VARCHAR(42),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crypto_tx_wallet ON crypto_transactions(wallet_id);
CREATE INDEX idx_crypto_tx_created ON crypto_transactions(created_at DESC);
