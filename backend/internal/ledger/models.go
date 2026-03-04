// Package ledger implements the double-entry general ledger, the financial
// core of TeslaPay's Core Banking System.
//
// All monetary values use shopspring/decimal to avoid floating-point errors.
// The ledger is append-only: corrections are made via reversal entries, never
// by modifying existing records.
package ledger

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AccountType represents the standard accounting classification.
type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"
)

// AccountCategory represents the business function of a ledger account.
type AccountCategory string

const (
	CategoryCustomerFunds   AccountCategory = "customer_funds"
	CategorySafeguardedFunds AccountCategory = "safeguarded_funds"
	CategoryFeeRevenue      AccountCategory = "fee_revenue"
	CategoryFXRevenue       AccountCategory = "fx_revenue"
	CategoryInterestExpense AccountCategory = "interest_expense"
	CategoryOperational     AccountCategory = "operational"
	CategorySettlement      AccountCategory = "settlement"
	CategorySuspense        AccountCategory = "suspense"
	CategoryControl         AccountCategory = "control"
)

// EntrySide indicates whether a ledger entry is a debit or credit.
type EntrySide string

const (
	Debit  EntrySide = "debit"
	Credit EntrySide = "credit"
)

// JournalEntryType categorizes the business purpose of a journal entry.
type JournalEntryType string

const (
	EntryTypePaymentDebit     JournalEntryType = "payment_debit"
	EntryTypePaymentCredit    JournalEntryType = "payment_credit"
	EntryTypeInternalTransfer JournalEntryType = "internal_transfer"
	EntryTypeCardAuth         JournalEntryType = "card_authorization"
	EntryTypeCardSettlement   JournalEntryType = "card_settlement"
	EntryTypeCardReversal     JournalEntryType = "card_reversal"
	EntryTypeFXExchange       JournalEntryType = "fx_exchange"
	EntryTypeCryptoBuy        JournalEntryType = "crypto_buy"
	EntryTypeCryptoSell       JournalEntryType = "crypto_sell"
	EntryTypeFee              JournalEntryType = "fee"
	EntryTypeInterest         JournalEntryType = "interest"
	EntryTypeAdjustment       JournalEntryType = "adjustment"
	EntryTypeReversal         JournalEntryType = "reversal"
	EntryTypeOpeningBalance   JournalEntryType = "opening_balance"
	EntryTypeClosing          JournalEntryType = "closing"
)

// JournalEntryStatus tracks the lifecycle of a journal entry.
type JournalEntryStatus string

const (
	StatusPending  JournalEntryStatus = "pending"
	StatusPosted   JournalEntryStatus = "posted"
	StatusReversed JournalEntryStatus = "reversed"
)

// ChartOfAccounts represents a single account in the general ledger's
// chart of accounts hierarchy.
type ChartOfAccounts struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Code      string          `json:"code" db:"code"`
	Name      string          `json:"name" db:"name"`
	Type      AccountType     `json:"type" db:"type"`
	Category  AccountCategory `json:"category" db:"category"`
	Currency  string          `json:"currency" db:"currency"`
	ParentID  *uuid.UUID      `json:"parent_id,omitempty" db:"parent_id"`
	IsSystem  bool            `json:"is_system" db:"is_system"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// JournalEntry represents the header for a balanced set of ledger entries.
// Every journal entry must have debits equal to credits in the same currency.
type JournalEntry struct {
	ID            uuid.UUID          `json:"id" db:"id"`
	PostingID     string             `json:"posting_id" db:"posting_id"` // Idempotency key.
	EffectiveDate time.Time          `json:"effective_date" db:"effective_date"`
	Description   string             `json:"description" db:"description"`
	EntryType     JournalEntryType   `json:"entry_type" db:"entry_type"`
	Status        JournalEntryStatus `json:"status" db:"status"`
	ReferenceType string             `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID   *uuid.UUID         `json:"reference_id,omitempty" db:"reference_id"`
	ReversalOf    *uuid.UUID         `json:"reversal_of,omitempty" db:"reversal_of"`
	Metadata      map[string]any     `json:"metadata,omitempty" db:"metadata"`
	CreatedBy     string             `json:"created_by" db:"created_by"`
	CreatedAt     time.Time          `json:"created_at" db:"created_at"`

	// Entries holds the individual debit/credit lines. Populated in memory
	// during posting and validation; not stored directly in journal_entries table.
	Entries []LedgerEntry `json:"entries,omitempty" db:"-"`
}

// LedgerEntry represents a single debit or credit line within a journal entry.
// This is the core double-entry record.
type LedgerEntry struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	JournalEntryID uuid.UUID       `json:"journal_entry_id" db:"journal_entry_id"`
	AccountID      uuid.UUID       `json:"account_id" db:"account_id"`
	EntrySide      EntrySide       `json:"entry_side" db:"entry_side"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`         // Always positive.
	Currency       string          `json:"currency" db:"currency"`
	BalanceAfter   decimal.Decimal `json:"balance_after" db:"balance_after"` // Running balance.
	SequenceNum    int64           `json:"sequence_num" db:"sequence_num"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// AccountBalance represents the materialized balance view (CQRS read model)
// for a single ledger account.
type AccountBalance struct {
	AccountID   uuid.UUID       `json:"account_id" db:"account_id"`
	Currency    string          `json:"currency" db:"currency"`
	Available   decimal.Decimal `json:"available" db:"available"`
	Pending     decimal.Decimal `json:"pending" db:"pending"`
	Reserved    decimal.Decimal `json:"reserved" db:"reserved"`
	Total       decimal.Decimal `json:"total" db:"total"`
	LastEntryID *uuid.UUID      `json:"last_entry_id,omitempty" db:"last_entry_id"`
	LastSequence int64          `json:"last_sequence" db:"last_sequence"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	Version     int64           `json:"version" db:"version"` // Optimistic locking.
}

// EventStoreEntry represents an immutable event in the ledger event store,
// used for event sourcing and audit trail.
type EventStoreEntry struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	AggregateID    uuid.UUID      `json:"aggregate_id" db:"aggregate_id"`
	AggregateType  string         `json:"aggregate_type" db:"aggregate_type"`
	EventType      string         `json:"event_type" db:"event_type"`
	EventData      map[string]any `json:"event_data" db:"event_data"`
	Metadata       map[string]any `json:"metadata,omitempty" db:"metadata"`
	SequenceNumber int64          `json:"sequence_number" db:"sequence_number"`
	Checksum       string         `json:"checksum" db:"checksum"` // SHA-256 chain.
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}

// PostTransactionRequest is the input for posting a new transaction to the ledger.
type PostTransactionRequest struct {
	PostingID     string           `json:"posting_id" validate:"required"`     // Idempotency key.
	EffectiveDate time.Time        `json:"effective_date" validate:"required"`
	Description   string           `json:"description" validate:"required,max=500"`
	EntryType     JournalEntryType `json:"entry_type" validate:"required"`
	ReferenceType string           `json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID       `json:"reference_id,omitempty"`
	CreatedBy     string           `json:"created_by" validate:"required"`
	Entries       []EntryLine      `json:"entries" validate:"required,min=2,dive"`
	Metadata      map[string]any   `json:"metadata,omitempty"`
}

// EntryLine is a single debit or credit line in a posting request.
type EntryLine struct {
	AccountID uuid.UUID       `json:"account_id" validate:"required"`
	Side      EntrySide       `json:"side" validate:"required,oneof=debit credit"`
	Amount    decimal.Decimal `json:"amount" validate:"required"`
	Currency  string          `json:"currency" validate:"required,len=3"`
}

// BalanceResponse is the API response for account balance queries.
type BalanceResponse struct {
	AccountID uuid.UUID       `json:"account_id"`
	Currency  string          `json:"currency"`
	Available string          `json:"available"`
	Pending   string          `json:"pending"`
	Total     string          `json:"total"`
}
