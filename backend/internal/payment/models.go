// Package payment implements the payment processing domain for TeslaPay,
// including internal transfers, SEPA payments, FX exchange, and scheduled payments.
package payment

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Payment types
const (
	PaymentTypeInternal = "internal"
	PaymentTypeSEPA     = "sepa"
)

// Payment statuses
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusCompleted  = "completed"
	PaymentStatusFailed     = "failed"
	PaymentStatusCancelled  = "cancelled"
)

// Schedule types
const (
	ScheduleTypeDaily   = "daily"
	ScheduleTypeWeekly  = "weekly"
	ScheduleTypeMonthly = "monthly"
)

// Payment is the core payment domain entity.
type Payment struct {
	ID                 uuid.UUID        `json:"id" db:"id"`
	SenderAccountID    uuid.UUID        `json:"sender_account_id" db:"sender_account_id"`
	RecipientAccountID *uuid.UUID       `json:"recipient_account_id,omitempty" db:"recipient_account_id"`
	RecipientIBAN      *string          `json:"recipient_iban,omitempty" db:"recipient_iban"`
	RecipientName      *string          `json:"recipient_name,omitempty" db:"recipient_name"`
	Amount             decimal.Decimal  `json:"amount" db:"amount"`
	Currency           string           `json:"currency" db:"currency"`
	Type               string           `json:"type" db:"type"`
	Status             string           `json:"status" db:"status"`
	Reference          *string          `json:"reference,omitempty" db:"reference"`
	Description        *string          `json:"description,omitempty" db:"description"`
	IdempotencyKey     *string          `json:"idempotency_key,omitempty" db:"idempotency_key"`
	FeeAmount          decimal.Decimal  `json:"fee_amount" db:"fee_amount"`
	FeeCurrency        string           `json:"fee_currency" db:"fee_currency"`
	FXRate             *decimal.Decimal `json:"fx_rate,omitempty" db:"fx_rate"`
	FXFromCurrency     *string          `json:"fx_from_currency,omitempty" db:"fx_from_currency"`
	FXToCurrency       *string          `json:"fx_to_currency,omitempty" db:"fx_to_currency"`
	CreatedAt          time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at" db:"updated_at"`
}

// ScheduledPayment is a recurring payment template.
type ScheduledPayment struct {
	ID                 uuid.UUID       `json:"id" db:"id"`
	AccountID          uuid.UUID       `json:"account_id" db:"account_id"`
	RecipientAccountID *uuid.UUID      `json:"recipient_account_id,omitempty" db:"recipient_account_id"`
	RecipientIBAN      *string         `json:"recipient_iban,omitempty" db:"recipient_iban"`
	RecipientName      *string         `json:"recipient_name,omitempty" db:"recipient_name"`
	Amount             decimal.Decimal `json:"amount" db:"amount"`
	Currency           string          `json:"currency" db:"currency"`
	Type               string          `json:"type" db:"type"`
	ScheduleType       string          `json:"schedule_type" db:"schedule_type"`
	Reference          *string         `json:"reference,omitempty" db:"reference"`
	Description        *string         `json:"description,omitempty" db:"description"`
	IsActive           bool            `json:"is_active" db:"is_active"`
	NextExecution      *time.Time      `json:"next_execution,omitempty" db:"next_execution"`
	LastExecution      *time.Time      `json:"last_execution,omitempty" db:"last_execution"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}

// FXQuote represents an FX exchange rate quote with expiry.
type FXQuote struct {
	ID              string          `json:"id"`
	FromCurrency    string          `json:"from_currency"`
	ToCurrency      string          `json:"to_currency"`
	Rate            decimal.Decimal `json:"rate"`
	Amount          decimal.Decimal `json:"amount"`
	ConvertedAmount decimal.Decimal `json:"converted_amount"`
	ExpiresAt       time.Time       `json:"expires_at"`
}

// -- Request types --

// InternalPaymentRequest is the request body for POST /payments/internal.
type InternalPaymentRequest struct {
	SenderAccountID    string `json:"sender_account_id" binding:"required"`
	RecipientAccountID string `json:"recipient_account_id" binding:"required"`
	Amount             string `json:"amount" binding:"required"`
	Currency           string `json:"currency" binding:"required,len=3"`
	Reference          string `json:"reference,omitempty"`
	Description        string `json:"description,omitempty"`
	IdempotencyKey     string `json:"idempotency_key,omitempty"`
}

// SEPAPaymentRequest is the request body for POST /payments/sepa.
type SEPAPaymentRequest struct {
	SenderAccountID string `json:"sender_account_id" binding:"required"`
	RecipientIBAN   string `json:"recipient_iban" binding:"required"`
	RecipientName   string `json:"recipient_name" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	Currency        string `json:"currency" binding:"required,len=3"`
	Reference       string `json:"reference,omitempty"`
	Description     string `json:"description,omitempty"`
	IdempotencyKey  string `json:"idempotency_key,omitempty"`
}

// FXQuoteRequest is the query for GET /payments/fx/quote.
type FXQuoteRequest struct {
	From   string `form:"from" binding:"required,len=3"`
	To     string `form:"to" binding:"required,len=3"`
	Amount string `form:"amount" binding:"required"`
}

// FXExecuteRequest is the request body for POST /payments/fx/execute.
type FXExecuteRequest struct {
	QuoteID   string `json:"quote_id" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
}

// ScheduledPaymentRequest is the request body for POST /payments/scheduled.
type ScheduledPaymentRequest struct {
	AccountID          string `json:"account_id" binding:"required"`
	RecipientAccountID string `json:"recipient_account_id,omitempty"`
	RecipientIBAN      string `json:"recipient_iban,omitempty"`
	RecipientName      string `json:"recipient_name,omitempty"`
	Amount             string `json:"amount" binding:"required"`
	Currency           string `json:"currency" binding:"required,len=3"`
	Type               string `json:"type" binding:"required,oneof=internal sepa"`
	ScheduleType       string `json:"schedule_type" binding:"required,oneof=daily weekly monthly"`
	Reference          string `json:"reference,omitempty"`
	Description        string `json:"description,omitempty"`
}

// -- Response types --

// PaymentResponse is the API response for payment operations.
type PaymentResponse struct {
	ID                 string `json:"id"`
	SenderAccountID    string `json:"sender_account_id"`
	RecipientAccountID string `json:"recipient_account_id,omitempty"`
	RecipientIBAN      string `json:"recipient_iban,omitempty"`
	RecipientName      string `json:"recipient_name,omitempty"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reference          string `json:"reference,omitempty"`
	Description        string `json:"description,omitempty"`
	FeeAmount          string `json:"fee_amount"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// toPaymentResponse converts a Payment domain object to API response.
func toPaymentResponse(p *Payment) *PaymentResponse {
	resp := &PaymentResponse{
		ID:              p.ID.String(),
		SenderAccountID: p.SenderAccountID.String(),
		Amount:          p.Amount.StringFixed(4),
		Currency:        p.Currency,
		Type:            p.Type,
		Status:          p.Status,
		FeeAmount:       p.FeeAmount.StringFixed(4),
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	}
	if p.RecipientAccountID != nil {
		resp.RecipientAccountID = p.RecipientAccountID.String()
	}
	if p.RecipientIBAN != nil {
		resp.RecipientIBAN = *p.RecipientIBAN
	}
	if p.RecipientName != nil {
		resp.RecipientName = *p.RecipientName
	}
	if p.Reference != nil {
		resp.Reference = *p.Reference
	}
	if p.Description != nil {
		resp.Description = *p.Description
	}
	return resp
}

// ListPaymentsResponse wraps a list of payments.
type ListPaymentsResponse struct {
	Data  []*PaymentResponse `json:"data"`
	Total int                `json:"total"`
}

// ScheduledPaymentResponse is the API response for scheduled payments.
type ScheduledPaymentResponse struct {
	ID                 string `json:"id"`
	AccountID          string `json:"account_id"`
	RecipientAccountID string `json:"recipient_account_id,omitempty"`
	RecipientIBAN      string `json:"recipient_iban,omitempty"`
	RecipientName      string `json:"recipient_name,omitempty"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	Type               string `json:"type"`
	ScheduleType       string `json:"schedule_type"`
	IsActive           bool   `json:"is_active"`
	NextExecution      string `json:"next_execution,omitempty"`
	CreatedAt          string `json:"created_at"`
}

func toScheduledPaymentResponse(sp *ScheduledPayment) *ScheduledPaymentResponse {
	resp := &ScheduledPaymentResponse{
		ID:           sp.ID.String(),
		AccountID:    sp.AccountID.String(),
		Amount:       sp.Amount.StringFixed(4),
		Currency:     sp.Currency,
		Type:         sp.Type,
		ScheduleType: sp.ScheduleType,
		IsActive:     sp.IsActive,
		CreatedAt:    sp.CreatedAt.Format(time.RFC3339),
	}
	if sp.RecipientAccountID != nil {
		resp.RecipientAccountID = sp.RecipientAccountID.String()
	}
	if sp.RecipientIBAN != nil {
		resp.RecipientIBAN = *sp.RecipientIBAN
	}
	if sp.RecipientName != nil {
		resp.RecipientName = *sp.RecipientName
	}
	if sp.NextExecution != nil {
		resp.NextExecution = sp.NextExecution.Format(time.RFC3339)
	}
	return resp
}
