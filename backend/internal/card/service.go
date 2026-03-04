package card

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/pkg/events"
	tpcrypto "github.com/teslapay/backend/pkg/crypto"
)

// devEncryptionKey is the AES-256 key for encrypting card numbers at rest.
// In production, this MUST be set via the CARD_ENCRYPTION_KEY env var (32 bytes).
// Default is a dev-only key; NEVER use in production.
var devEncryptionKey = []byte("teslapay-card-key-dev-32bytekey!") // exactly 32 bytes

func getEncryptionKey() []byte {
	if k := os.Getenv("CARD_ENCRYPTION_KEY"); len(k) == 32 {
		return []byte(k)
	}
	return devEncryptionKey
}

// Service implements the card business logic.
type Service struct {
	repo     *Repository
	producer *events.Producer
	logger   *zap.Logger
}

// NewService creates a new card service.
func NewService(repo *Repository, producer *events.Producer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

// IssueVirtualCard issues a new virtual Mastercard.
func (s *Service) IssueVirtualCard(ctx context.Context, req *IssueCardRequest) (*Card, error) {
	return s.issueCard(ctx, req, CardTypeVirtual, CardStatusActive)
}

// IssuePhysicalCard issues a new physical Mastercard (starts as pending_delivery).
func (s *Service) IssuePhysicalCard(ctx context.Context, req *IssueCardRequest) (*Card, error) {
	return s.issueCard(ctx, req, CardTypePhysical, CardStatusPendingDelivery)
}

// issueCard is the shared card issuance logic.
func (s *Service) issueCard(ctx context.Context, req *IssueCardRequest, cardType, initialStatus string) (*Card, error) {
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, common.NewValidationError("Invalid account_id", "")
	}

	// Generate card number with Mastercard BIN 5425.
	cardNumber, err := tpcrypto.GenerateCardNumber("5425")
	if err != nil {
		return nil, fmt.Errorf("generate card number: %w", err)
	}

	// Validate Luhn.
	if !tpcrypto.ValidateLuhn(cardNumber) {
		return nil, fmt.Errorf("generated card number failed Luhn check")
	}

	// Encrypt card number for storage.
	encryptedNumber, err := tpcrypto.EncryptCardNumber(cardNumber, getEncryptionKey())
	if err != nil {
		return nil, fmt.Errorf("encrypt card number: %w", err)
	}

	lastFour := cardNumber[len(cardNumber)-4:]

	// Generate CVV and hash it — never store plaintext CVV.
	cvv, err := tpcrypto.GenerateCVV()
	if err != nil {
		return nil, fmt.Errorf("generate CVV: %w", err)
	}
	cvvHash := tpcrypto.HashCVV(cvv)

	// Expiry: 3 years from now.
	now := time.Now().UTC()
	expiryDate := now.AddDate(3, 0, 0)

	var subAccountID *uuid.UUID
	if req.SubAccountID != "" {
		sid, err := uuid.Parse(req.SubAccountID)
		if err != nil {
			return nil, common.NewValidationError("Invalid sub_account_id", "")
		}
		subAccountID = &sid
	}

	card := &Card{
		ID:                  uuid.New(),
		AccountID:           accountID,
		SubAccountID:        subAccountID,
		CardNumberEncrypted: encryptedNumber,
		LastFour:            lastFour,
		ExpiryMonth:         int(expiryDate.Month()),
		ExpiryYear:          expiryDate.Year(),
		CVVHash:             cvvHash,
		CardholderName:      req.CardholderName,
		Type:                cardType,
		Status:              initialStatus,
		DailyLimit:          decimal.NewFromInt(5000),
		MonthlyLimit:        decimal.NewFromInt(25000),
		DailySpent:          decimal.Zero,
		MonthlySpent:        decimal.Zero,
		IsContactless:       true,
		IsOnline:            true,
		IsATM:               true,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.CreateCard(ctx, card); err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("card issued",
			zap.String("card_id", card.ID.String()),
			zap.String("account_id", accountID.String()),
			zap.String("type", cardType),
			zap.String("last_four", lastFour),
		)
	}

	s.publishCardEvent(ctx, card, "card.issued")

	return card, nil
}

// ActivateCard activates a card by verifying the last 4 digits.
func (s *Service) ActivateCard(ctx context.Context, cardID uuid.UUID, req *ActivateCardRequest) (*Card, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}

	if card.LastFour != req.LastFour {
		return nil, common.NewValidationError("Card last four digits do not match", "")
	}

	if card.Status != CardStatusPendingDelivery {
		return nil, common.NewBusinessError(
			"CARD_002",
			"Card Cannot Be Activated",
			fmt.Sprintf("Card with status '%s' cannot be activated", card.Status),
		)
	}

	if err := s.repo.UpdateCardStatus(ctx, cardID, CardStatusActive); err != nil {
		return nil, fmt.Errorf("activate card: %w", err)
	}

	card.Status = CardStatusActive
	s.publishCardEvent(ctx, card, "card.activated")

	return card, nil
}

// FreezeCard temporarily freezes an active card.
func (s *Service) FreezeCard(ctx context.Context, cardID uuid.UUID) (*Card, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}

	if card.Status != CardStatusActive {
		return nil, common.NewBusinessError(
			common.ErrCodeCardFrozen,
			"Card Cannot Be Frozen",
			fmt.Sprintf("Only active cards can be frozen; current status: '%s'", card.Status),
		)
	}

	if err := s.repo.UpdateCardStatus(ctx, cardID, CardStatusFrozen); err != nil {
		return nil, fmt.Errorf("freeze card: %w", err)
	}

	card.Status = CardStatusFrozen
	s.publishCardEvent(ctx, card, "card.frozen")

	return card, nil
}

// UnfreezeCard unfreezes a frozen card.
func (s *Service) UnfreezeCard(ctx context.Context, cardID uuid.UUID) (*Card, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}

	if card.Status != CardStatusFrozen {
		return nil, common.NewBusinessError(
			"CARD_003",
			"Card Cannot Be Unfrozen",
			fmt.Sprintf("Only frozen cards can be unfrozen; current status: '%s'", card.Status),
		)
	}

	if err := s.repo.UpdateCardStatus(ctx, cardID, CardStatusActive); err != nil {
		return nil, fmt.Errorf("unfreeze card: %w", err)
	}

	card.Status = CardStatusActive
	s.publishCardEvent(ctx, card, "card.unfrozen")

	return card, nil
}

// BlockCard permanently blocks a card (irreversible).
func (s *Service) BlockCard(ctx context.Context, cardID uuid.UUID) (*Card, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}

	if card.Status == CardStatusBlocked {
		return nil, common.NewConflictError("Card is already blocked")
	}
	if card.Status == CardStatusExpired || card.Status == CardStatusCancelled {
		return nil, common.NewBusinessError(
			"CARD_004",
			"Card Cannot Be Blocked",
			fmt.Sprintf("Card with status '%s' cannot be blocked", card.Status),
		)
	}

	if err := s.repo.UpdateCardStatus(ctx, cardID, CardStatusBlocked); err != nil {
		return nil, fmt.Errorf("block card: %w", err)
	}

	card.Status = CardStatusBlocked
	s.publishCardEvent(ctx, card, "card.blocked")

	return card, nil
}

// UpdateControls updates card spending limits and permissions.
func (s *Service) UpdateControls(ctx context.Context, cardID uuid.UUID, req *CardControlsRequest) (*Card, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}

	// Apply updates to the existing card, with defaults from current values.
	dailyLimit := card.DailyLimit
	monthlyLimit := card.MonthlyLimit
	isContactless := card.IsContactless
	isOnline := card.IsOnline
	isATM := card.IsATM
	allowedCountries := card.AllowedCountries
	blockedMCC := card.BlockedMCCCodes

	if req.DailyLimit != nil {
		dl, err := decimal.NewFromString(*req.DailyLimit)
		if err != nil || dl.LessThan(decimal.Zero) {
			return nil, common.NewValidationError("Invalid daily_limit: must be a non-negative decimal", "")
		}
		dailyLimit = dl
	}
	if req.MonthlyLimit != nil {
		ml, err := decimal.NewFromString(*req.MonthlyLimit)
		if err != nil || ml.LessThan(decimal.Zero) {
			return nil, common.NewValidationError("Invalid monthly_limit: must be a non-negative decimal", "")
		}
		monthlyLimit = ml
	}
	if req.IsContactless != nil {
		isContactless = *req.IsContactless
	}
	if req.IsOnline != nil {
		isOnline = *req.IsOnline
	}
	if req.IsATM != nil {
		isATM = *req.IsATM
	}
	if req.AllowedCountries != nil {
		allowedCountries = req.AllowedCountries
	}
	if req.BlockedMCCCodes != nil {
		blockedMCC = req.BlockedMCCCodes
	}

	update := &CardControlsUpdate{
		DailyLimit:       dailyLimit,
		MonthlyLimit:     monthlyLimit,
		IsContactless:    isContactless,
		IsOnline:         isOnline,
		IsATM:            isATM,
		AllowedCountries: allowedCountries,
		BlockedMCCCodes:  blockedMCC,
	}

	if err := s.repo.UpdateCardControls(ctx, cardID, update); err != nil {
		return nil, fmt.Errorf("update card controls: %w", err)
	}

	// Update in-memory card for response.
	card.DailyLimit = dailyLimit
	card.MonthlyLimit = monthlyLimit
	card.IsContactless = isContactless
	card.IsOnline = isOnline
	card.IsATM = isATM
	card.AllowedCountries = allowedCountries
	card.BlockedMCCCodes = blockedMCC

	return card, nil
}

// GetCard retrieves a card by ID (masked response — never exposes full number).
func (s *Service) GetCard(ctx context.Context, cardID uuid.UUID) (*Card, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}
	return card, nil
}

// ListCards returns all cards for an account.
func (s *Service) ListCards(ctx context.Context, accountID uuid.UUID) ([]*Card, error) {
	return s.repo.GetCardsByAccountID(ctx, accountID)
}

// GetCardTransactions returns paginated transactions for a card.
func (s *Service) GetCardTransactions(ctx context.Context, cardID uuid.UUID, limit, offset int) ([]*CardTransaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	// Verify card exists.
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	if card == nil {
		return nil, common.NewNotFoundError("Card not found", "")
	}
	return s.repo.GetCardTransactions(ctx, cardID, limit, offset)
}

// publishCardEvent publishes a card domain event to Kafka (best-effort).
func (s *Service) publishCardEvent(ctx context.Context, c *Card, eventType string) {
	if s.producer == nil {
		return
	}

	eventData := map[string]any{
		"card_id":    c.ID.String(),
		"account_id": c.AccountID.String(),
		"type":       c.Type,
		"status":     c.Status,
		"last_four":  c.LastFour,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to marshal card event data", zap.Error(err))
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
		Source:        "card-service",
		AggregateID:   c.ID.String(),
		AggregateType: "card",
		Data:          data,
		Metadata: events.EventMetadata{
			CorrelationID: correlationID,
			ActorID:       c.AccountID.String(),
			ActorType:     "account",
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := s.producer.Publish(ctx, events.TopicCardEvents, event); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to publish card event",
				zap.String("card_id", c.ID.String()),
				zap.String("event_type", eventType),
				zap.Error(err),
			)
		}
	}
}
