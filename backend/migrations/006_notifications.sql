CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT,
    data JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notif_user ON notifications(user_id);
CREATE INDEX idx_notif_created ON notifications(created_at DESC);
CREATE INDEX idx_notif_unread ON notifications(user_id, read_at) WHERE read_at IS NULL;

CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id UUID PRIMARY KEY,
    push_enabled BOOLEAN DEFAULT true,
    email_enabled BOOLEAN DEFAULT true,
    sms_enabled BOOLEAN DEFAULT false,
    payment_alerts BOOLEAN DEFAULT true,
    card_alerts BOOLEAN DEFAULT true,
    kyc_alerts BOOLEAN DEFAULT true,
    marketing BOOLEAN DEFAULT false,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
