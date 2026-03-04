package ledger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/pkg/crypto"
	"github.com/teslapay/backend/pkg/events"
)

// Service encapsulates the business logic for the double-entry ledger.
type Service struct {
	repo     *Repository
	producer *events.Producer
	logger   *zap.Logger
}

// NewService creates a new ledger service.
func NewService(repo *Repository, producer *events.Producer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

// PostTransaction creates a balanced journal entry with its ledger entries.
// This is the primary write path for all financial operations in TeslaPay.
//
// Guarantees:
//   - Idempotent: duplicate PostingIDs return the existing journal entry.
//   - Balanced: total debits must equal total credits per currency.
//   - Atomic: all-or-nothing via database transaction.
//   - Auditable: an event is published to Kafka after successful posting.
func (s *Service) PostTransaction(ctx context.Context, req *PostTransactionRequest) (*JournalEntry, error) {
	// 1. Validate the request.
	if err := s.validatePostRequest(req); err != nil {
		return nil, err
	}

	// 2. Check idempotency: if this posting_id already exists, return existing.
	existing, err := s.repo.GetJournalEntryByPostingID(ctx, req.PostingID)
	if err != nil {
		return nil, fmt.Errorf("idempotency check: %w", err)
	}
	if existing != nil {
		s.logger.Info("idempotent duplicate detected",
			zap.String("posting_id", req.PostingID),
			zap.String("journal_id", existing.ID.String()),
		)
		return existing, nil
	}

	// 3. Validate all accounts exist and are active.
	for _, entry := range req.Entries {
		acct, err := s.repo.GetAccountByID(ctx, entry.AccountID)
		if err != nil {
			return nil, common.NewNotFoundError(
				fmt.Sprintf("Ledger account %s not found", entry.AccountID), "",
			)
		}
		if !acct.IsActive {
			return nil, common.NewBusinessError(
				"LEDGER_001", "Account Inactive",
				fmt.Sprintf("Ledger account %s is inactive", entry.AccountID),
			)
		}
		if acct.Currency != entry.Currency {
			return nil, common.NewValidationError(
				fmt.Sprintf("Currency mismatch: account %s is %s but entry specifies %s",
					entry.AccountID, acct.Currency, entry.Currency),
				"",
			)
		}
	}

	// 4. For debit entries on liability accounts (customer funds), verify sufficient balance.
	for _, entry := range req.Entries {
		if entry.Side == Debit {
			bal, err := s.repo.GetBalance(ctx, entry.AccountID)
			if err != nil {
				// Balance row might not exist yet for new accounts.
				s.logger.Debug("no balance row found for account",
					zap.String("account_id", entry.AccountID.String()),
				)
				continue
			}
			if bal.Available.LessThan(entry.Amount) {
				return nil, common.NewBusinessError(
					common.ErrCodeInsufficientFunds,
					"Insufficient Funds",
					fmt.Sprintf("Account %s balance is %s, but operation requires %s",
						entry.AccountID, bal.Available.StringFixed(2), entry.Amount.StringFixed(2)),
				)
			}
		}
	}

	// 5. Build the journal entry and ledger entries.
	now := time.Now().UTC()
	je := &JournalEntry{
		ID:            uuid.New(),
		PostingID:     req.PostingID,
		EffectiveDate: req.EffectiveDate,
		Description:   req.Description,
		EntryType:     req.EntryType,
		Status:        StatusPosted,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Metadata:      req.Metadata,
		CreatedBy:     req.CreatedBy,
		CreatedAt:     now,
	}

	ledgerEntries := make([]LedgerEntry, len(req.Entries))
	for i, line := range req.Entries {
		ledgerEntries[i] = LedgerEntry{
			ID:        uuid.New(),
			AccountID: line.AccountID,
			EntrySide: line.Side,
			Amount:    line.Amount,
			Currency:  line.Currency,
		}
	}

	// 6. Post atomically to the database.
	if err := s.repo.PostJournalEntry(ctx, je, ledgerEntries); err != nil {
		return nil, fmt.Errorf("post journal entry: %w", err)
	}

	je.Entries = ledgerEntries

	// 7. Publish event to Kafka (best-effort; the DB is the source of truth).
	s.publishEntryPostedEvent(ctx, je)

	return je, nil
}

// GetBalance retrieves the current balance for a ledger account.
func (s *Service) GetBalance(ctx context.Context, accountID uuid.UUID) (*BalanceResponse, error) {
	bal, err := s.repo.GetBalance(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return &BalanceResponse{
		AccountID: bal.AccountID,
		Currency:  bal.Currency,
		Available: bal.Available.StringFixed(4),
		Pending:   bal.Pending.StringFixed(4),
		Total:     bal.Total.StringFixed(4),
	}, nil
}

// GetTransactionHistory returns paginated ledger entries for an account.
func (s *Service) GetTransactionHistory(
	ctx context.Context, accountID uuid.UUID, limit int, afterSequence int64,
) ([]LedgerEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.GetLedgerEntriesForAccount(ctx, accountID, limit, afterSequence)
}

// validatePostRequest performs business rule validation on a posting request.
func (s *Service) validatePostRequest(req *PostTransactionRequest) error {
	if req.PostingID == "" {
		return common.NewValidationError("posting_id is required", "")
	}
	if len(req.Entries) < 2 {
		return common.NewValidationError(
			"A journal entry must have at least 2 entry lines (debit and credit)", "",
		)
	}

	// Validate all amounts are positive.
	for i, entry := range req.Entries {
		if entry.Amount.LessThanOrEqual(decimal.Zero) {
			return common.NewValidationError(
				fmt.Sprintf("Entry %d: amount must be positive, got %s", i, entry.Amount.String()),
				"",
			)
		}
	}

	// Validate double-entry: total debits must equal total credits per currency.
	balanceByCurrency := make(map[string]decimal.Decimal)
	for _, entry := range req.Entries {
		if entry.Side == Debit {
			balanceByCurrency[entry.Currency] = balanceByCurrency[entry.Currency].Add(entry.Amount)
		} else if entry.Side == Credit {
			balanceByCurrency[entry.Currency] = balanceByCurrency[entry.Currency].Sub(entry.Amount)
		} else {
			return common.NewValidationError(
				fmt.Sprintf("Invalid entry side: %s (must be 'debit' or 'credit')", entry.Side), "",
			)
		}
	}

	for currency, balance := range balanceByCurrency {
		if !balance.IsZero() {
			return common.NewValidationError(
				fmt.Sprintf("Double-entry violation: %s entries are unbalanced by %s",
					currency, balance.Abs().String()),
				"",
			)
		}
	}

	return nil
}

// ValidateDoubleEntry checks that a given set of entry lines satisfies the
// double-entry accounting invariant. This is a pure validation function
// that does not touch the database.
func (s *Service) ValidateDoubleEntry(entries []EntryLine) error {
	if len(entries) < 2 {
		return fmt.Errorf("at least 2 entries required, got %d", len(entries))
	}

	balanceByCurrency := make(map[string]decimal.Decimal)
	for _, entry := range entries {
		if entry.Amount.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("entry amount must be positive: %s", entry.Amount.String())
		}
		switch entry.Side {
		case Debit:
			balanceByCurrency[entry.Currency] = balanceByCurrency[entry.Currency].Add(entry.Amount)
		case Credit:
			balanceByCurrency[entry.Currency] = balanceByCurrency[entry.Currency].Sub(entry.Amount)
		default:
			return fmt.Errorf("invalid entry side: %s", entry.Side)
		}
	}

	for currency, balance := range balanceByCurrency {
		if !balance.IsZero() {
			return fmt.Errorf("%s entries unbalanced by %s", currency, balance.Abs().String())
		}
	}

	return nil
}

// publishEntryPostedEvent sends a ledger event to Kafka for downstream consumers.
func (s *Service) publishEntryPostedEvent(ctx context.Context, je *JournalEntry) {
	if s.producer == nil {
		return
	}

	eventData := map[string]any{
		"journal_entry_id": je.ID.String(),
		"posting_id":       je.PostingID,
		"entry_type":       je.EntryType,
		"effective_date":   je.EffectiveDate.Format(time.DateOnly),
		"description":      je.Description,
		"entry_count":      len(je.Entries),
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		s.logger.Error("failed to marshal event data", zap.Error(err))
		return
	}

	event := &events.Event{
		ID:            uuid.New().String(),
		Type:          "entry.posted",
		Source:        "ledger-service",
		AggregateID:   je.ID.String(),
		AggregateType: "journal",
		Data:          data,
		Metadata: events.EventMetadata{
			CorrelationID: common.RequestIDFromContext(ctx),
			ActorID:       je.CreatedBy,
			ActorType:     "service",
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := s.producer.Publish(ctx, events.TopicLedgerEvents, event); err != nil {
		// Log but do not fail: the database is the source of truth.
		s.logger.Error("failed to publish ledger event",
			zap.String("journal_id", je.ID.String()),
			zap.Error(err),
		)
	}
}

// ComputeEventChecksum generates a SHA-256 checksum for event sourcing chain integrity.
func ComputeEventChecksum(eventData map[string]any, previousChecksum string) string {
	data, _ := json.Marshal(eventData)
	return crypto.SHA256Hash(string(data) + previousChecksum)
}
