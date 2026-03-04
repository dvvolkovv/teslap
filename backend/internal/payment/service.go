package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/pkg/events"
)

// fxRates contains hardcoded exchange rates (EUR as base currency).
var fxRates = map[string]map[string]decimal.Decimal{
	"EUR": {
		"USD": decimal.NewFromFloat(1.08),
		"GBP": decimal.NewFromFloat(0.86),
		"PLN": decimal.NewFromFloat(4.30),
		"CHF": decimal.NewFromFloat(0.96),
		"EUR": decimal.NewFromFloat(1.00),
	},
	"USD": {
		"EUR": decimal.NewFromFloat(0.9259),
		"GBP": decimal.NewFromFloat(0.7963),
		"USD": decimal.NewFromFloat(1.00),
	},
	"GBP": {
		"EUR": decimal.NewFromFloat(1.1628),
		"USD": decimal.NewFromFloat(1.2563),
		"GBP": decimal.NewFromFloat(1.00),
	},
}

// Service implements the payment business logic.
type Service struct {
	repo     *Repository
	producer *events.Producer
	logger   *zap.Logger
}

// NewService creates a new payment service.
func NewService(repo *Repository, producer *events.Producer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

// CreateInternalPayment processes a transfer between two TeslaPay accounts.
func (s *Service) CreateInternalPayment(ctx context.Context, req *InternalPaymentRequest) (*Payment, error) {
	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, common.NewValidationError("Invalid amount: must be a positive decimal number", "")
	}

	// Parse account IDs
	senderID, err := uuid.Parse(req.SenderAccountID)
	if err != nil {
		return nil, common.NewValidationError("Invalid sender_account_id", "")
	}
	recipientID, err := uuid.Parse(req.RecipientAccountID)
	if err != nil {
		return nil, common.NewValidationError("Invalid recipient_account_id", "")
	}

	// Sender and recipient must be different
	if senderID == recipientID {
		return nil, common.NewValidationError("Sender and recipient accounts must be different", "")
	}

	// Check idempotency
	if req.IdempotencyKey != "" {
		existing, err := s.repo.GetPaymentByIdempotencyKey(ctx, req.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("idempotency check: %w", err)
		}
		if existing != nil {
			if s.logger != nil {
				s.logger.Info("idempotent payment request detected",
					zap.String("idempotency_key", req.IdempotencyKey),
					zap.String("payment_id", existing.ID.String()),
				)
			}
			return existing, nil
		}
	}

	now := time.Now().UTC()
	idempotencyKey := stringOrNil(req.IdempotencyKey)
	reference := stringOrNil(req.Reference)
	description := stringOrNil(req.Description)

	payment := &Payment{
		ID:                 uuid.New(),
		SenderAccountID:    senderID,
		RecipientAccountID: &recipientID,
		Amount:             amount,
		Currency:           req.Currency,
		Type:               PaymentTypeInternal,
		Status:             PaymentStatusCompleted, // Internal transfers complete immediately
		Reference:          reference,
		Description:        description,
		IdempotencyKey:     idempotencyKey,
		FeeAmount:          decimal.Zero,
		FeeCurrency:        req.Currency,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("internal payment created",
			zap.String("payment_id", payment.ID.String()),
			zap.String("sender", senderID.String()),
			zap.String("recipient", recipientID.String()),
			zap.String("amount", amount.StringFixed(4)),
			zap.String("currency", req.Currency),
		)
	}

	// Publish event (best-effort)
	s.publishPaymentEvent(ctx, payment, "payment.completed")

	return payment, nil
}

// CreateSEPAPayment initiates an external SEPA payment.
func (s *Service) CreateSEPAPayment(ctx context.Context, req *SEPAPaymentRequest) (*Payment, error) {
	// Validate IBAN format (basic check: length and alphanumeric after country code)
	if err := validateIBAN(req.RecipientIBAN); err != nil {
		return nil, err
	}

	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, common.NewValidationError("Invalid amount: must be a positive decimal number", "")
	}

	// Parse sender account ID
	senderID, err := uuid.Parse(req.SenderAccountID)
	if err != nil {
		return nil, common.NewValidationError("Invalid sender_account_id", "")
	}

	// Check idempotency
	if req.IdempotencyKey != "" {
		existing, err := s.repo.GetPaymentByIdempotencyKey(ctx, req.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("idempotency check: %w", err)
		}
		if existing != nil {
			return existing, nil
		}
	}

	now := time.Now().UTC()
	iban := req.RecipientIBAN
	name := req.RecipientName
	reference := stringOrNil(req.Reference)
	description := stringOrNil(req.Description)
	idempotencyKey := stringOrNil(req.IdempotencyKey)

	payment := &Payment{
		ID:              uuid.New(),
		SenderAccountID: senderID,
		RecipientIBAN:   &iban,
		RecipientName:   &name,
		Amount:          amount,
		Currency:        req.Currency,
		Type:            PaymentTypeSEPA,
		Status:          PaymentStatusPending,
		Reference:       reference,
		Description:     description,
		IdempotencyKey:  idempotencyKey,
		FeeAmount:       decimal.Zero,
		FeeCurrency:     req.Currency,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("create sepa payment: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("SEPA payment initiated",
			zap.String("payment_id", payment.ID.String()),
			zap.String("sender", senderID.String()),
			zap.String("iban", req.RecipientIBAN),
			zap.String("amount", amount.StringFixed(4)),
		)
	}

	// Publish payment.initiated event
	s.publishPaymentEvent(ctx, payment, "payment.initiated")

	return payment, nil
}

// GetFXQuote returns an FX exchange rate quote valid for 30 seconds.
func (s *Service) GetFXQuote(ctx context.Context, from, to, amountStr string) (*FXQuote, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	amount, err := decimal.NewFromString(amountStr)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, common.NewValidationError("Invalid amount: must be a positive decimal number", "")
	}

	rates, ok := fxRates[from]
	if !ok {
		return nil, common.NewValidationError(
			fmt.Sprintf("Currency %s is not supported for FX conversion", from), "",
		)
	}
	rate, ok := rates[to]
	if !ok {
		return nil, common.NewValidationError(
			fmt.Sprintf("FX pair %s/%s is not supported", from, to), "",
		)
	}

	convertedAmount := amount.Mul(rate)

	return &FXQuote{
		ID:              uuid.New().String(),
		FromCurrency:    from,
		ToCurrency:      to,
		Rate:            rate,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
		ExpiresAt:       time.Now().UTC().Add(30 * time.Second),
	}, nil
}

// ExecuteFX executes an FX conversion. In Phase 1 this is a placeholder.
func (s *Service) ExecuteFX(ctx context.Context, req *FXExecuteRequest) error {
	// In production: validate quote is not expired, execute conversion via ledger
	if s.logger != nil {
		s.logger.Info("FX execution requested (placeholder)",
			zap.String("quote_id", req.QuoteID),
			zap.String("account_id", req.AccountID),
		)
	}
	return nil
}

// GetPayment retrieves a single payment by ID.
func (s *Service) GetPayment(ctx context.Context, id uuid.UUID) (*Payment, error) {
	payment, err := s.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}
	if payment == nil {
		return nil, common.NewNotFoundError("Payment not found", "")
	}
	return payment, nil
}

// ListPayments returns paginated payments for an account.
func (s *Service) ListPayments(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*Payment, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetPaymentsByAccountID(ctx, accountID, limit, offset)
}

// CreateScheduledPayment creates a new scheduled payment.
func (s *Service) CreateScheduledPayment(ctx context.Context, req *ScheduledPaymentRequest) (*ScheduledPayment, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, common.NewValidationError("Invalid amount: must be a positive decimal number", "")
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, common.NewValidationError("Invalid account_id", "")
	}

	now := time.Now().UTC()
	nextExecution := computeNextExecution(req.ScheduleType, now)
	reference := stringOrNil(req.Reference)
	description := stringOrNil(req.Description)

	sp := &ScheduledPayment{
		ID:            uuid.New(),
		AccountID:     accountID,
		Amount:        amount,
		Currency:      req.Currency,
		Type:          req.Type,
		ScheduleType:  req.ScheduleType,
		Reference:     reference,
		Description:   description,
		IsActive:      true,
		NextExecution: &nextExecution,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if req.RecipientAccountID != "" {
		rid, err := uuid.Parse(req.RecipientAccountID)
		if err != nil {
			return nil, common.NewValidationError("Invalid recipient_account_id", "")
		}
		sp.RecipientAccountID = &rid
	}
	if req.RecipientIBAN != "" {
		iban := req.RecipientIBAN
		sp.RecipientIBAN = &iban
	}
	if req.RecipientName != "" {
		name := req.RecipientName
		sp.RecipientName = &name
	}

	if err := s.repo.CreateScheduledPayment(ctx, sp); err != nil {
		return nil, fmt.Errorf("create scheduled payment: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("scheduled payment created",
			zap.String("scheduled_id", sp.ID.String()),
			zap.String("account_id", accountID.String()),
			zap.String("schedule_type", req.ScheduleType),
		)
	}

	return sp, nil
}

// ListScheduledPayments returns all scheduled payments for an account.
func (s *Service) ListScheduledPayments(ctx context.Context, accountID uuid.UUID) ([]*ScheduledPayment, error) {
	return s.repo.GetScheduledPayments(ctx, accountID)
}

// publishPaymentEvent publishes a payment event to Kafka (best-effort).
func (s *Service) publishPaymentEvent(ctx context.Context, p *Payment, eventType string) {
	if s.producer == nil {
		return
	}

	eventData := map[string]any{
		"payment_id":        p.ID.String(),
		"sender_account_id": p.SenderAccountID.String(),
		"amount":            p.Amount.StringFixed(4),
		"currency":          p.Currency,
		"type":              p.Type,
		"status":            p.Status,
	}
	if p.RecipientAccountID != nil {
		eventData["recipient_account_id"] = p.RecipientAccountID.String()
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to marshal payment event data", zap.Error(err))
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
		Source:        "payment-service",
		AggregateID:   p.ID.String(),
		AggregateType: "payment",
		Data:          data,
		Metadata: events.EventMetadata{
			CorrelationID: correlationID,
			ActorID:       p.SenderAccountID.String(),
			ActorType:     "account",
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := s.producer.Publish(ctx, events.TopicPaymentEvents, event); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to publish payment event",
				zap.String("payment_id", p.ID.String()),
				zap.String("event_type", eventType),
				zap.Error(err),
			)
		}
	}
}

// validateIBAN performs basic IBAN format validation.
// Full validation would use a dedicated IBAN library in production.
func validateIBAN(iban string) error {
	iban = strings.ToUpper(strings.ReplaceAll(iban, " ", ""))
	if len(iban) < 15 || len(iban) > 34 {
		return common.NewBusinessError(
			common.ErrCodeInvalidIBAN,
			"Invalid IBAN",
			"IBAN must be between 15 and 34 characters",
		)
	}
	// Country code: first 2 chars must be letters
	if iban[0] < 'A' || iban[0] > 'Z' || iban[1] < 'A' || iban[1] > 'Z' {
		return common.NewBusinessError(
			common.ErrCodeInvalidIBAN,
			"Invalid IBAN",
			"IBAN must start with a valid country code",
		)
	}
	// Check digits: next 2 chars must be digits
	if iban[2] < '0' || iban[2] > '9' || iban[3] < '0' || iban[3] > '9' {
		return common.NewBusinessError(
			common.ErrCodeInvalidIBAN,
			"Invalid IBAN",
			"IBAN check digits must be numeric",
		)
	}
	return nil
}

// computeNextExecution computes the next execution time for a scheduled payment.
func computeNextExecution(scheduleType string, from time.Time) time.Time {
	switch scheduleType {
	case ScheduleTypeDaily:
		return from.Add(24 * time.Hour)
	case ScheduleTypeWeekly:
		return from.Add(7 * 24 * time.Hour)
	case ScheduleTypeMonthly:
		return from.AddDate(0, 1, 0)
	default:
		return from.Add(24 * time.Hour)
	}
}

// stringOrNil returns a pointer to the string if non-empty, or nil.
func stringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
