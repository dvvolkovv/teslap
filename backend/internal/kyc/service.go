package kyc

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

// Service implements the KYC business logic.
type Service struct {
	repo     *Repository
	producer *events.Producer
	logger   *zap.Logger
}

// NewService creates a new KYC service.
func NewService(repo *Repository, producer *events.Producer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

// StartVerification initiates a new KYC verification for a user.
func (s *Service) StartVerification(ctx context.Context, req *KYCStartRequest) (*KYCStartResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id format", "")
	}

	existing, err := s.repo.GetKYCByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check existing kyc: %w", err)
	}
	if existing != nil && (existing.Status == StatusPending || existing.Status == StatusInReview || existing.Status == StatusApproved) {
		return nil, common.NewConflictError("KYC verification already active or approved")
	}

	applicantID := "sumsub_" + uuid.New().String()
	sdkToken := "token_" + uuid.New().String()

	now := time.Now().UTC()
	record := &KYCRecord{
		ID:          uuid.New(),
		UserID:      userID,
		Provider:    "sumsub",
		ApplicantID: applicantID,
		Level:       req.Level,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateKYCRecord(ctx, record); err != nil {
		return nil, fmt.Errorf("create kyc record: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("kyc verification started",
			zap.String("user_id", userID.String()),
			zap.String("applicant_id", applicantID),
			zap.String("level", req.Level),
		)
	}

	return &KYCStartResponse{
		ApplicantID: applicantID,
		SDKToken:    sdkToken,
		FlowURL:     "https://api.sumsub.com/idensic/msdk/applicant/" + applicantID,
	}, nil
}

// GetStatus returns the current KYC status for a user.
func (s *Service) GetStatus(ctx context.Context, userID uuid.UUID) (*KYCStatusResponse, error) {
	record, err := s.repo.GetKYCByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get kyc status: %w", err)
	}
	if record == nil {
		return nil, common.NewNotFoundError("KYC record not found", "")
	}

	resp := &KYCStatusResponse{
		Level:      record.Level,
		Status:     record.Status,
		VerifiedAt: record.VerifiedAt,
	}
	if len(record.ReviewResult) > 0 {
		resp.ReviewResult = record.ReviewResult
	}
	return resp, nil
}

// HandleWebhook processes an incoming Sumsub webhook payload.
func (s *Service) HandleWebhook(ctx context.Context, payload *KYCWebhookPayload) error {
	// Placeholder signature validation — in production, verify HMAC signature from headers.
	if s.logger != nil {
		s.logger.Info("kyc webhook received",
			zap.String("applicant_id", payload.ApplicantID),
			zap.String("review_status", payload.ReviewStatus),
		)
	}

	record, err := s.repo.GetKYCByApplicantID(ctx, payload.ApplicantID)
	if err != nil {
		return fmt.Errorf("lookup kyc by applicant_id: %w", err)
	}
	if record == nil {
		return common.NewNotFoundError("KYC record not found for applicant", "")
	}

	var newStatus string
	switch payload.ReviewStatus {
	case "completed":
		newStatus = StatusApproved
	case "declined":
		newStatus = StatusRejected
	default:
		newStatus = StatusInReview
	}

	if err := s.repo.UpdateKYCStatus(ctx, payload.ApplicantID, newStatus, payload.ReviewResult); err != nil {
		return fmt.Errorf("update kyc status: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("kyc status updated",
			zap.String("applicant_id", payload.ApplicantID),
			zap.String("new_status", newStatus),
		)
	}

	// Update local record status for event publishing.
	record.Status = newStatus

	switch newStatus {
	case StatusApproved:
		s.publishKYCEvent(ctx, record, "kyc.approved")
	case StatusRejected:
		s.publishKYCEvent(ctx, record, "kyc.rejected")
	}

	return nil
}

// publishKYCEvent publishes a KYC event to Kafka (best-effort).
func (s *Service) publishKYCEvent(ctx context.Context, record *KYCRecord, eventType string) {
	if s.producer == nil {
		return
	}

	eventData := map[string]any{
		"user_id":      record.UserID.String(),
		"applicant_id": record.ApplicantID,
		"level":        record.Level,
		"status":       record.Status,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to marshal kyc event data", zap.Error(err))
		}
		return
	}

	correlationID := ""
	if ctx != nil {
		correlationID = common.RequestIDFromContext(ctx)
	}

	event := &events.Event{
		ID:            uuid.New().String(),
		Type:          eventType,
		Source:        "kyc-service",
		AggregateID:   record.UserID.String(),
		AggregateType: "user",
		Data:          data,
		Metadata: events.EventMetadata{
			CorrelationID: correlationID,
			ActorID:       record.UserID.String(),
			ActorType:     "user",
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := s.producer.Publish(ctx, events.TopicKYCEvents, event); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to publish kyc event",
				zap.String("user_id", record.UserID.String()),
				zap.String("event_type", eventType),
				zap.Error(err),
			)
		}
	}
}
