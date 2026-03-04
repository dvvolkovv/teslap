package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/pkg/events"
)

// Service implements the notification business logic.
type Service struct {
	repo     *Repository
	producer *events.Producer
	logger   *zap.Logger
}

// NewService creates a new notification service.
func NewService(repo *Repository, producer *events.Producer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

// SendNotification creates and dispatches a notification for a user.
func (s *Service) SendNotification(ctx context.Context, req *SendNotificationRequest) (*Notification, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	// Fetch user preferences to determine whether to deliver on the requested channel.
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}

	// Determine channel — skip delivery but still record in_app if push is disabled.
	channel := req.Type
	if prefs != nil {
		switch req.Type {
		case TypePush:
			if !prefs.PushEnabled {
				// Push is disabled: fall back to in_app record so the user still sees the message.
				channel = TypeInApp
			}
		case TypeEmail:
			if !prefs.EmailEnabled {
				// Email disabled: skip entirely.
				if s.logger != nil {
					s.logger.Info("notification skipped: email disabled",
						zap.String("user_id", userID.String()),
					)
				}
				return nil, common.NewBusinessError(
					"NOTIF_001",
					"Notification Skipped",
					"Email notifications are disabled for this user",
				)
			}
		}
	}

	now := time.Now().UTC()
	n := &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      req.Type,
		Channel:   channel,
		Title:     req.Title,
		Body:      req.Body,
		Data:      req.Data,
		Status:    StatusPending,
		SentAt:    &now,
		CreatedAt: now,
	}

	if err := s.repo.CreateNotification(ctx, n); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	// In production: dispatch via FCM/APNs/SendGrid/Twilio (placeholder — log only).
	if s.logger != nil {
		s.logger.Info("notification created",
			zap.String("notification_id", n.ID.String()),
			zap.String("user_id", userID.String()),
			zap.String("type", n.Type),
			zap.String("channel", n.Channel),
			zap.String("title", n.Title),
		)
	}

	s.publishNotificationEvent(ctx, n.ID.String(), userID.String(), "notification.sent", map[string]any{
		"notification_id": n.ID.String(),
		"user_id":         userID.String(),
		"type":            n.Type,
		"channel":         n.Channel,
		"title":           n.Title,
	})

	return n, nil
}

// ListNotifications returns paginated notifications for a user.
func (s *Service) ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetNotificationsByUserID(ctx, userID, limit, offset)
}

// MarkAsRead marks a notification as read.
func (s *Service) MarkAsRead(ctx context.Context, notifID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, notifID)
}

// GetPreferences returns notification preferences for a user.
// If no preferences record exists, default preferences are returned.
func (s *Service) GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	if prefs == nil {
		return defaultPreferences(userID), nil
	}
	return prefs, nil
}

// UpdatePreferences applies partial updates to a user's notification preferences.
func (s *Service) UpdatePreferences(ctx context.Context, userID uuid.UUID, req *PreferencesUpdateRequest) (*NotificationPreferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	if prefs == nil {
		prefs = defaultPreferences(userID)
	}

	// Apply only the fields that were explicitly provided.
	if req.PushEnabled != nil {
		prefs.PushEnabled = *req.PushEnabled
	}
	if req.EmailEnabled != nil {
		prefs.EmailEnabled = *req.EmailEnabled
	}
	if req.SMSEnabled != nil {
		prefs.SMSEnabled = *req.SMSEnabled
	}
	if req.PaymentAlerts != nil {
		prefs.PaymentAlerts = *req.PaymentAlerts
	}
	if req.CardAlerts != nil {
		prefs.CardAlerts = *req.CardAlerts
	}
	if req.KYCAlerts != nil {
		prefs.KYCAlerts = *req.KYCAlerts
	}
	if req.Marketing != nil {
		prefs.Marketing = *req.Marketing
	}

	prefs.UserID = userID
	prefs.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpsertPreferences(ctx, prefs); err != nil {
		return nil, fmt.Errorf("upsert preferences: %w", err)
	}

	return prefs, nil
}

// defaultPreferences returns the system-default notification preferences for a new user.
func defaultPreferences(userID uuid.UUID) *NotificationPreferences {
	return &NotificationPreferences{
		UserID:        userID,
		PushEnabled:   true,
		EmailEnabled:  true,
		SMSEnabled:    false,
		PaymentAlerts: true,
		CardAlerts:    true,
		KYCAlerts:     true,
		Marketing:     false,
		UpdatedAt:     time.Now().UTC(),
	}
}

// publishNotificationEvent publishes a notification domain event to Kafka (best-effort).
func (s *Service) publishNotificationEvent(ctx context.Context, aggregateID, actorID, eventType string, eventData map[string]any) {
	if s.producer == nil {
		return
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to marshal notification event data", zap.Error(err))
		}
		return
	}

	correlationID := ""
	if ctx != nil {
		correlationID = common.RequestIDFromContext(ctx)
	}

	event := events.Event{
		ID:            uuid.New().String(),
		Type:          eventType,
		Source:        "notification-service",
		AggregateID:   aggregateID,
		AggregateType: "notification",
		Data:          data,
		Metadata: events.EventMetadata{
			CorrelationID: correlationID,
			ActorID:       actorID,
		},
	}

	if err := s.producer.Publish(ctx, events.TopicNotificationCommands, &event); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to publish notification event",
				zap.String("event_type", eventType),
				zap.Error(err),
			)
		}
	}
}
