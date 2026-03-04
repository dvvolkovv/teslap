// Package card implements the card management domain for TeslaPay,
// including virtual/physical card issuance, lifecycle management,
// spending controls, and transaction tracking.
package card

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Card types
const (
	CardTypeVirtual  = "virtual"
	CardTypePhysical = "physical"
)

// Card statuses
const (
	CardStatusActive          = "active"
	CardStatusFrozen          = "frozen"
	CardStatusBlocked         = "blocked"
	CardStatusExpired         = "expired"
	CardStatusCancelled       = "cancelled"
	CardStatusPendingDelivery = "pending_delivery"
)

// Card transaction types
const (
	TxTypePurchase = "purchase"
	TxTypeATM      = "atm"
	TxTypeRefund   = "refund"
	TxTypeReversal = "reversal"
)

// Card transaction statuses
const (
	TxStatusPending   = "pending"
	TxStatusCompleted = "completed"
	TxStatusDeclined  = "declined"
	TxStatusReversed  = "reversed"
)

// Card is the core card domain entity.
type Card struct {
	ID                  uuid.UUID       `json:"id" db:"id"`
	AccountID           uuid.UUID       `json:"account_id" db:"account_id"`
	SubAccountID        *uuid.UUID      `json:"sub_account_id,omitempty" db:"sub_account_id"`
	CardNumberEncrypted string          `json:"-" db:"card_number_encrypted"` // NEVER in JSON
	LastFour            string          `json:"last_four" db:"last_four"`
	ExpiryMonth         int             `json:"expiry_month" db:"expiry_month"`
	ExpiryYear          int             `json:"expiry_year" db:"expiry_year"`
	CVVHash             string          `json:"-" db:"cvv_hash"` // NEVER in JSON
	CardholderName      string          `json:"cardholder_name" db:"cardholder_name"`
	Type                string          `json:"type" db:"type"`
	Status              string          `json:"status" db:"status"`
	DailyLimit          decimal.Decimal `json:"daily_limit" db:"daily_limit"`
	MonthlyLimit        decimal.Decimal `json:"monthly_limit" db:"monthly_limit"`
	DailySpent          decimal.Decimal `json:"daily_spent" db:"daily_spent"`
	MonthlySpent        decimal.Decimal `json:"monthly_spent" db:"monthly_spent"`
	IsContactless       bool            `json:"is_contactless" db:"is_contactless"`
	IsOnline            bool            `json:"is_online" db:"is_online"`
	IsATM               bool            `json:"is_atm" db:"is_atm"`
	AllowedCountries    []string        `json:"allowed_countries,omitempty" db:"allowed_countries"`
	BlockedMCCCodes     []int           `json:"blocked_mcc_codes,omitempty" db:"blocked_mcc_codes"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}

// CardTransaction represents a card transaction record.
type CardTransaction struct {
	ID                uuid.UUID        `json:"id" db:"id"`
	CardID            uuid.UUID        `json:"card_id" db:"card_id"`
	MerchantName      *string          `json:"merchant_name,omitempty" db:"merchant_name"`
	MerchantMCC       *int             `json:"merchant_mcc,omitempty" db:"merchant_mcc"`
	Amount            decimal.Decimal  `json:"amount" db:"amount"`
	Currency          string           `json:"currency" db:"currency"`
	OriginalAmount    *decimal.Decimal `json:"original_amount,omitempty" db:"original_amount"`
	OriginalCurrency  *string          `json:"original_currency,omitempty" db:"original_currency"`
	Status            string           `json:"status" db:"status"`
	Type              string           `json:"type" db:"type"`
	Country           *string          `json:"country,omitempty" db:"country"`
	AuthorizationCode *string          `json:"authorization_code,omitempty" db:"authorization_code"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
}

// CardControlsUpdate holds parsed, validated values for updating card controls.
type CardControlsUpdate struct {
	DailyLimit       decimal.Decimal
	MonthlyLimit     decimal.Decimal
	IsContactless    bool
	IsOnline         bool
	IsATM            bool
	AllowedCountries []string
	BlockedMCCCodes  []int
}

// -- Request types --

// IssueCardRequest is the request body for POST /cards/virtual or /cards/physical.
type IssueCardRequest struct {
	AccountID      string `json:"account_id" binding:"required"`
	SubAccountID   string `json:"sub_account_id,omitempty"`
	CardholderName string `json:"cardholder_name" binding:"required"`
}

// ActivateCardRequest is the request body for POST /cards/:card_id/activate.
type ActivateCardRequest struct {
	LastFour string `json:"last_four" binding:"required,len=4"`
}

// CardControlsRequest is the request body for PUT /cards/:card_id/controls.
type CardControlsRequest struct {
	DailyLimit       *string  `json:"daily_limit,omitempty"`
	MonthlyLimit     *string  `json:"monthly_limit,omitempty"`
	IsContactless    *bool    `json:"is_contactless,omitempty"`
	IsOnline         *bool    `json:"is_online,omitempty"`
	IsATM            *bool    `json:"is_atm,omitempty"`
	AllowedCountries []string `json:"allowed_countries,omitempty"`
	BlockedMCCCodes  []int    `json:"blocked_mcc_codes,omitempty"`
}

// -- Response types --

// CardResponse is the safe API response — never exposes full card number or CVV.
type CardResponse struct {
	ID             string `json:"id"`
	AccountID      string `json:"account_id"`
	SubAccountID   string `json:"sub_account_id,omitempty"`
	MaskedNumber   string `json:"masked_number"`
	LastFour       string `json:"last_four"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CardholderName string `json:"cardholder_name"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	DailyLimit     string `json:"daily_limit"`
	MonthlyLimit   string `json:"monthly_limit"`
	DailySpent     string `json:"daily_spent"`
	MonthlySpent   string `json:"monthly_spent"`
	IsContactless  bool   `json:"is_contactless"`
	IsOnline       bool   `json:"is_online"`
	IsATM          bool   `json:"is_atm"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ListCardsResponse wraps a list of cards.
type ListCardsResponse struct {
	Data  []*CardResponse `json:"data"`
	Total int             `json:"total"`
}

// CardTransactionResponse is the API response for a card transaction.
type CardTransactionResponse struct {
	ID               string `json:"id"`
	CardID           string `json:"card_id"`
	MerchantName     string `json:"merchant_name,omitempty"`
	MerchantMCC      int    `json:"merchant_mcc,omitempty"`
	Amount           string `json:"amount"`
	Currency         string `json:"currency"`
	OriginalAmount   string `json:"original_amount,omitempty"`
	OriginalCurrency string `json:"original_currency,omitempty"`
	Status           string `json:"status"`
	Type             string `json:"type"`
	Country          string `json:"country,omitempty"`
	CreatedAt        string `json:"created_at"`
}

// ListCardTransactionsResponse wraps a list of card transactions.
type ListCardTransactionsResponse struct {
	Data  []*CardTransactionResponse `json:"data"`
	Total int                        `json:"total"`
}

// toCardResponse converts a Card domain object to API response.
// The masked_number is computed from the BIN "5425" + last_four since we don't
// store the full card number in memory outside of issuance.
func toCardResponse(c *Card) *CardResponse {
	resp := &CardResponse{
		ID:             c.ID.String(),
		AccountID:      c.AccountID.String(),
		MaskedNumber:   fmt.Sprintf("5425 **** **** %s", c.LastFour),
		LastFour:       c.LastFour,
		ExpiryMonth:    c.ExpiryMonth,
		ExpiryYear:     c.ExpiryYear,
		CardholderName: c.CardholderName,
		Type:           c.Type,
		Status:         c.Status,
		DailyLimit:     c.DailyLimit.StringFixed(4),
		MonthlyLimit:   c.MonthlyLimit.StringFixed(4),
		DailySpent:     c.DailySpent.StringFixed(4),
		MonthlySpent:   c.MonthlySpent.StringFixed(4),
		IsContactless:  c.IsContactless,
		IsOnline:       c.IsOnline,
		IsATM:          c.IsATM,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      c.UpdatedAt.Format(time.RFC3339),
	}
	if c.SubAccountID != nil {
		resp.SubAccountID = c.SubAccountID.String()
	}
	return resp
}

// toCardTransactionResponse converts a CardTransaction to API response.
func toCardTransactionResponse(tx *CardTransaction) *CardTransactionResponse {
	resp := &CardTransactionResponse{
		ID:        tx.ID.String(),
		CardID:    tx.CardID.String(),
		Amount:    tx.Amount.StringFixed(4),
		Currency:  tx.Currency,
		Status:    tx.Status,
		Type:      tx.Type,
		CreatedAt: tx.CreatedAt.Format(time.RFC3339),
	}
	if tx.MerchantName != nil {
		resp.MerchantName = *tx.MerchantName
	}
	if tx.MerchantMCC != nil {
		resp.MerchantMCC = *tx.MerchantMCC
	}
	if tx.OriginalAmount != nil {
		resp.OriginalAmount = tx.OriginalAmount.StringFixed(4)
	}
	if tx.OriginalCurrency != nil {
		resp.OriginalCurrency = *tx.OriginalCurrency
	}
	if tx.Country != nil {
		resp.Country = *tx.Country
	}
	return resp
}
