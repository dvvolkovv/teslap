package notification

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the notification service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new notification repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreateNotification inserts a new notification record.
func (r *Repository) CreateNotification(ctx context.Context, n *Notification) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO notifications (id, user_id, type, channel, title, body, data, status, sent_at, read_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		n.ID, n.UserID, n.Type, n.Channel, n.Title, n.Body, n.Data, n.Status, n.SentAt, n.ReadAt, n.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	return nil
}

// GetNotificationsByUserID returns paginated notifications for a user ordered by creation date descending.
func (r *Repository) GetNotificationsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Notification, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, user_id, type, channel, title, body, data, status, sent_at, read_at, created_at
		FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query notifications by user: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Channel, &n.Title, &n.Body, &n.Data, &n.Status, &n.SentAt, &n.ReadAt, &n.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, &n)
	}
	return notifications, rows.Err()
}

// MarkAsRead sets the read_at timestamp on a notification.
func (r *Repository) MarkAsRead(ctx context.Context, notifID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE notifications SET read_at = NOW() WHERE id = $1
	`, notifID)
	if err != nil {
		return fmt.Errorf("mark notification as read: %w", err)
	}
	return nil
}

// GetPreferences fetches notification preferences for a user.
// Returns nil, nil if no preferences record exists.
func (r *Repository) GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error) {
	var prefs NotificationPreferences
	err := r.db.Pool.QueryRow(ctx, `
		SELECT user_id, push_enabled, email_enabled, sms_enabled, payment_alerts, card_alerts, kyc_alerts, marketing, updated_at
		FROM notification_preferences WHERE user_id = $1
	`, userID).Scan(
		&prefs.UserID, &prefs.PushEnabled, &prefs.EmailEnabled, &prefs.SMSEnabled,
		&prefs.PaymentAlerts, &prefs.CardAlerts, &prefs.KYCAlerts, &prefs.Marketing, &prefs.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query notification preferences: %w", err)
	}
	return &prefs, nil
}

// UpsertPreferences inserts or updates notification preferences for a user.
func (r *Repository) UpsertPreferences(ctx context.Context, prefs *NotificationPreferences) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO notification_preferences (user_id, push_enabled, email_enabled, sms_enabled, payment_alerts, card_alerts, kyc_alerts, marketing, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
		  push_enabled = EXCLUDED.push_enabled,
		  email_enabled = EXCLUDED.email_enabled,
		  sms_enabled = EXCLUDED.sms_enabled,
		  payment_alerts = EXCLUDED.payment_alerts,
		  card_alerts = EXCLUDED.card_alerts,
		  kyc_alerts = EXCLUDED.kyc_alerts,
		  marketing = EXCLUDED.marketing,
		  updated_at = EXCLUDED.updated_at
	`,
		prefs.UserID, prefs.PushEnabled, prefs.EmailEnabled, prefs.SMSEnabled,
		prefs.PaymentAlerts, prefs.CardAlerts, prefs.KYCAlerts, prefs.Marketing,
		prefs.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert notification preferences: %w", err)
	}
	return nil
}

