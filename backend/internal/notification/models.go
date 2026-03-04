package notification

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Notification represents a single notification record.
type Notification struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Type      string          `json:"type" db:"type"`
	Channel   string          `json:"channel" db:"channel"`
	Title     string          `json:"title" db:"title"`
	Body      string          `json:"body,omitempty" db:"body"`
	Data      json.RawMessage `json:"data,omitempty" db:"data"`
	Status    string          `json:"status" db:"status"`
	SentAt    *time.Time      `json:"sent_at,omitempty" db:"sent_at"`
	ReadAt    *time.Time      `json:"read_at,omitempty" db:"read_at"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// NotificationPreferences stores per-user notification settings.
type NotificationPreferences struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	PushEnabled   bool      `json:"push_enabled" db:"push_enabled"`
	EmailEnabled  bool      `json:"email_enabled" db:"email_enabled"`
	SMSEnabled    bool      `json:"sms_enabled" db:"sms_enabled"`
	PaymentAlerts bool      `json:"payment_alerts" db:"payment_alerts"`
	CardAlerts    bool      `json:"card_alerts" db:"card_alerts"`
	KYCAlerts     bool      `json:"kyc_alerts" db:"kyc_alerts"`
	Marketing     bool      `json:"marketing" db:"marketing"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Notification type constants.
const (
	TypePush  = "push"
	TypeEmail = "email"
	TypeSMS   = "sms"
	TypeInApp = "in_app"
)

// Notification status constants.
const (
	StatusPending   = "pending"
	StatusSent      = "sent"
	StatusDelivered = "delivered"
	StatusFailed    = "failed"
)

// -- Request / Response types --

// SendNotificationRequest matches POST (internal only).
type SendNotificationRequest struct {
	UserID string          `json:"user_id" binding:"required"`
	Type   string          `json:"type" binding:"required,oneof=push email sms in_app"`
	Title  string          `json:"title" binding:"required"`
	Body   string          `json:"body"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// PreferencesUpdateRequest matches PUT /notifications/preferences.
type PreferencesUpdateRequest struct {
	PushEnabled   *bool `json:"push_enabled,omitempty"`
	EmailEnabled  *bool `json:"email_enabled,omitempty"`
	SMSEnabled    *bool `json:"sms_enabled,omitempty"`
	PaymentAlerts *bool `json:"payment_alerts,omitempty"`
	CardAlerts    *bool `json:"card_alerts,omitempty"`
	KYCAlerts     *bool `json:"kyc_alerts,omitempty"`
	Marketing     *bool `json:"marketing,omitempty"`
}

// ListNotificationsResponse wraps the notification list.
type ListNotificationsResponse struct {
	Data  []*Notification `json:"data"`
	Total int             `json:"total"`
}
