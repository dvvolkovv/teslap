CREATE TABLE IF NOT EXISTS kyc_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider VARCHAR(50) NOT NULL DEFAULT 'sumsub',
    applicant_id VARCHAR(100) UNIQUE,
    level VARCHAR(20) NOT NULL DEFAULT 'basic',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    review_result JSONB,
    rejection_reasons TEXT[],
    documents_provided TEXT[],
    verified_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_kyc_user ON kyc_records(user_id);
CREATE INDEX idx_kyc_applicant ON kyc_records(applicant_id);
CREATE INDEX idx_kyc_status ON kyc_records(status);
