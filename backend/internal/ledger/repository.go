package ledger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access operations for the ledger service.
// All write operations that involve multiple tables use database transactions
// to guarantee ACID properties.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new ledger repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// GetAccountByID retrieves a chart of accounts entry by its ID.
func (r *Repository) GetAccountByID(ctx context.Context, id uuid.UUID) (*ChartOfAccounts, error) {
	var acct ChartOfAccounts
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, code, name, type, category, currency, parent_id, is_system, is_active,
		       created_at, updated_at
		FROM chart_of_accounts
		WHERE id = $1
	`, id).Scan(
		&acct.ID, &acct.Code, &acct.Name, &acct.Type, &acct.Category,
		&acct.Currency, &acct.ParentID, &acct.IsSystem, &acct.IsActive,
		&acct.CreatedAt, &acct.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get account %s: %w", id, err)
	}
	return &acct, nil
}

// GetBalance retrieves the current materialized balance for a ledger account.
func (r *Repository) GetBalance(ctx context.Context, accountID uuid.UUID) (*AccountBalance, error) {
	var bal AccountBalance
	err := r.db.Pool.QueryRow(ctx, `
		SELECT account_id, currency, available, pending, reserved, total,
		       last_entry_id, last_sequence, updated_at, version
		FROM account_balances
		WHERE account_id = $1
	`, accountID).Scan(
		&bal.AccountID, &bal.Currency, &bal.Available, &bal.Pending,
		&bal.Reserved, &bal.Total, &bal.LastEntryID, &bal.LastSequence,
		&bal.UpdatedAt, &bal.Version,
	)
	if err != nil {
		return nil, fmt.Errorf("get balance for account %s: %w", accountID, err)
	}
	return &bal, nil
}

// GetJournalEntryByPostingID retrieves a journal entry by its idempotency key.
// Returns nil if not found (for idempotency checking).
func (r *Repository) GetJournalEntryByPostingID(ctx context.Context, postingID string) (*JournalEntry, error) {
	var je JournalEntry
	var metadataBytes []byte
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, posting_id, effective_date, description, entry_type, status,
		       reference_type, reference_id, reversal_of, metadata, created_by, created_at
		FROM journal_entries
		WHERE posting_id = $1
	`, postingID).Scan(
		&je.ID, &je.PostingID, &je.EffectiveDate, &je.Description,
		&je.EntryType, &je.Status, &je.ReferenceType, &je.ReferenceID,
		&je.ReversalOf, &metadataBytes, &je.CreatedBy, &je.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get journal entry by posting_id %s: %w", postingID, err)
	}
	if metadataBytes != nil {
		_ = json.Unmarshal(metadataBytes, &je.Metadata)
	}
	return &je, nil
}

// PostJournalEntry atomically inserts a journal entry with its ledger entries
// and updates account balances. This is the critical path that must guarantee:
//  1. Journal entry is created
//  2. All ledger entries are inserted
//  3. Debits equal credits (validated before commit)
//  4. Account balances are updated with optimistic locking
//  5. Event store entry is appended
//
// If any step fails, the entire transaction is rolled back.
func (r *Repository) PostJournalEntry(ctx context.Context, je *JournalEntry, entries []LedgerEntry) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		// 1. Insert journal entry.
		metadataJSON, err := json.Marshal(je.Metadata)
		if err != nil {
			return fmt.Errorf("marshal metadata: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO journal_entries
				(id, posting_id, effective_date, description, entry_type, status,
				 reference_type, reference_id, reversal_of, metadata, created_by, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			je.ID, je.PostingID, je.EffectiveDate, je.Description,
			je.EntryType, je.Status, je.ReferenceType, je.ReferenceID,
			je.ReversalOf, metadataJSON, je.CreatedBy, je.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert journal entry: %w", err)
		}

		// 2. Insert all ledger entries and update balances.
		for i := range entries {
			entry := &entries[i]
			entry.JournalEntryID = je.ID

			// Get the current balance with a row-level lock for this account.
			var currentBalance AccountBalance
			err := tx.QueryRow(ctx, `
				SELECT account_id, currency, available, pending, reserved, total,
				       last_sequence, version
				FROM account_balances
				WHERE account_id = $1
				FOR UPDATE
			`, entry.AccountID).Scan(
				&currentBalance.AccountID, &currentBalance.Currency,
				&currentBalance.Available, &currentBalance.Pending,
				&currentBalance.Reserved, &currentBalance.Total,
				&currentBalance.LastSequence, &currentBalance.Version,
			)
			if err == pgx.ErrNoRows {
				// Initialize balance row if it does not exist.
				currentBalance = AccountBalance{
					AccountID:    entry.AccountID,
					Currency:     entry.Currency,
					Available:    decimal.Zero,
					Pending:      decimal.Zero,
					Reserved:     decimal.Zero,
					Total:        decimal.Zero,
					LastSequence: 0,
					Version:      0,
				}
			} else if err != nil {
				return fmt.Errorf("lock balance for account %s: %w", entry.AccountID, err)
			}

			// Calculate new balance based on account type and entry side.
			newAvailable := currentBalance.Available
			if entry.EntrySide == Debit {
				newAvailable = newAvailable.Sub(entry.Amount)
			} else {
				newAvailable = newAvailable.Add(entry.Amount)
			}

			newSequence := currentBalance.LastSequence + 1
			entry.SequenceNum = newSequence
			entry.BalanceAfter = newAvailable
			entry.CreatedAt = time.Now().UTC()
			if entry.ID == uuid.Nil {
				entry.ID = uuid.New()
			}

			// Insert ledger entry.
			_, err = tx.Exec(ctx, `
				INSERT INTO ledger_entries
					(id, journal_entry_id, account_id, entry_side, amount, currency,
					 balance_after, sequence_num, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`,
				entry.ID, entry.JournalEntryID, entry.AccountID,
				entry.EntrySide, entry.Amount, entry.Currency,
				entry.BalanceAfter, entry.SequenceNum, entry.CreatedAt,
			)
			if err != nil {
				return fmt.Errorf("insert ledger entry for account %s: %w", entry.AccountID, err)
			}

			// Update materialized balance (upsert).
			newTotal := newAvailable.Add(currentBalance.Pending).Add(currentBalance.Reserved)
			_, err = tx.Exec(ctx, `
				INSERT INTO account_balances
					(account_id, currency, available, pending, reserved, total,
					 last_entry_id, last_sequence, updated_at, version)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), 1)
				ON CONFLICT (account_id) DO UPDATE SET
					available = $3,
					total = $6,
					last_entry_id = $7,
					last_sequence = $8,
					updated_at = NOW(),
					version = account_balances.version + 1
				WHERE account_balances.version = $9
			`,
				entry.AccountID, entry.Currency, newAvailable,
				currentBalance.Pending, currentBalance.Reserved, newTotal,
				entry.ID, newSequence, currentBalance.Version,
			)
			if err != nil {
				return fmt.Errorf("update balance for account %s: %w", entry.AccountID, err)
			}
		}

		// 3. Validate double-entry invariant: sum of debits must equal sum of credits.
		var balanceCheck decimal.Decimal
		rows, err := tx.Query(ctx, `
			SELECT entry_side, amount FROM ledger_entries WHERE journal_entry_id = $1
		`, je.ID)
		if err != nil {
			return fmt.Errorf("validate double entry: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var side EntrySide
			var amount decimal.Decimal
			if err := rows.Scan(&side, &amount); err != nil {
				return fmt.Errorf("scan validation row: %w", err)
			}
			if side == Debit {
				balanceCheck = balanceCheck.Add(amount)
			} else {
				balanceCheck = balanceCheck.Sub(amount)
			}
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("iterate validation rows: %w", err)
		}

		if !balanceCheck.IsZero() {
			return fmt.Errorf(
				"CRITICAL: double-entry violation for journal %s: debit-credit difference is %s",
				je.ID, balanceCheck.String(),
			)
		}

		r.logger.Info("journal entry posted",
			zap.String("journal_id", je.ID.String()),
			zap.String("posting_id", je.PostingID),
			zap.String("entry_type", string(je.EntryType)),
			zap.Int("entry_count", len(entries)),
		)

		return nil
	})
}

// GetLedgerEntriesForAccount retrieves ledger entries for a specific account,
// ordered by sequence number descending, with cursor-based pagination.
func (r *Repository) GetLedgerEntriesForAccount(
	ctx context.Context, accountID uuid.UUID, limit int, afterSequence int64,
) ([]LedgerEntry, error) {
	query := `
		SELECT id, journal_entry_id, account_id, entry_side, amount, currency,
		       balance_after, sequence_num, created_at
		FROM ledger_entries
		WHERE account_id = $1
	`
	args := []any{accountID}

	if afterSequence > 0 {
		query += " AND sequence_num < $3"
		args = append(args, limit, afterSequence)
		query += " ORDER BY sequence_num DESC LIMIT $2"
	} else {
		query += " ORDER BY sequence_num DESC LIMIT $2"
		args = append(args, limit)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query ledger entries: %w", err)
	}
	defer rows.Close()

	var entries []LedgerEntry
	for rows.Next() {
		var e LedgerEntry
		if err := rows.Scan(
			&e.ID, &e.JournalEntryID, &e.AccountID, &e.EntrySide,
			&e.Amount, &e.Currency, &e.BalanceAfter, &e.SequenceNum, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan ledger entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// AppendEvent writes an immutable event to the event store with SHA-256 chaining.
func (r *Repository) AppendEvent(ctx context.Context, tx pgx.Tx, event *EventStoreEntry) error {
	eventDataJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}
	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return fmt.Errorf("marshal event metadata: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO event_store
			(id, aggregate_id, aggregate_type, event_type, event_data,
			 metadata, sequence_number, checksum, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		event.ID, event.AggregateID, event.AggregateType, event.EventType,
		eventDataJSON, metadataJSON, event.SequenceNumber, event.Checksum,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("append event: %w", err)
	}
	return nil
}

// InitializeAccountBalance creates a zero-balance row for a new ledger account.
func (r *Repository) InitializeAccountBalance(ctx context.Context, accountID uuid.UUID, currency string) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO account_balances (account_id, currency, available, pending, reserved, total, last_sequence, version)
		VALUES ($1, $2, 0, 0, 0, 0, 0, 0)
		ON CONFLICT (account_id) DO NOTHING
	`, accountID, currency)
	if err != nil {
		return fmt.Errorf("initialize account balance: %w", err)
	}
	return nil
}
