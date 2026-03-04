// Package account implements the account management domain for TeslaPay,
// including user profiles, multi-currency sub-accounts, IBAN generation,
// and account tiers.
package account

import (
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// User represents a TeslaPay user profile in the account database.
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ExternalID   string     `json:"external_id" db:"external_id"`
	TierID       uuid.UUID  `json:"tier_id" db:"tier_id"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	DateOfBirth  time.Time  `json:"date_of_birth" db:"date_of_birth"`
	Nationality  string     `json:"nationality,omitempty" db:"nationality"`
	TaxResidency string     `json:"tax_residency,omitempty" db:"tax_residency"`
	AddressLine1 string     `json:"address_line1,omitempty" db:"address_line1"`
	AddressLine2 string     `json:"address_line2,omitempty" db:"address_line2"`
	City         string     `json:"city,omitempty" db:"city"`
	PostalCode   string     `json:"postal_code,omitempty" db:"postal_code"`
	Country      string     `json:"country" db:"country"`
	Language     string     `json:"language" db:"language"`
	KYCStatus    string     `json:"kyc_status" db:"kyc_status"`
	KYCLevel     int        `json:"kyc_level" db:"kyc_level"`
	RiskScore    int        `json:"risk_score" db:"risk_score"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Account is a container for sub-accounts, one per user.
type Account struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	AccountNumber string     `json:"account_number" db:"account_number"`
	Status        string     `json:"status" db:"status"`
	OpenedAt      time.Time  `json:"opened_at" db:"opened_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty" db:"closed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// SubAccount represents a single-currency account with an optional IBAN.
type SubAccount struct {
	ID              uuid.UUID `json:"id" db:"id"`
	AccountID       uuid.UUID `json:"account_id" db:"account_id"`
	Currency        string    `json:"currency" db:"currency"`
	IBAN            *string   `json:"iban" db:"iban"`
	BIC             *string   `json:"bic,omitempty" db:"bic"`
	LedgerAccountID uuid.UUID `json:"ledger_account_id" db:"ledger_account_id"`
	Status          string    `json:"status" db:"status"`
	IsDefault       bool      `json:"is_default" db:"is_default"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// AccountTier defines the feature limits for an account tier level.
type AccountTier struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	Name                 string          `json:"name" db:"name"`
	MonthlyFee           decimal.Decimal `json:"monthly_fee" db:"monthly_fee"`
	DailyTransferLimit   decimal.Decimal `json:"daily_transfer_limit" db:"daily_transfer_limit"`
	MonthlyTransferLimit decimal.Decimal `json:"monthly_transfer_limit" db:"monthly_transfer_limit"`
	DailyCardLimit       decimal.Decimal `json:"daily_card_limit" db:"daily_card_limit"`
	MonthlyATMLimit      decimal.Decimal `json:"monthly_atm_limit" db:"monthly_atm_limit"`
	FreeATMWithdrawals   int             `json:"free_atm_withdrawals" db:"free_atm_withdrawals"`
	FXMarkupPercent      decimal.Decimal `json:"fx_markup_percent" db:"fx_markup_percent"`
	MaxSubAccounts       int             `json:"max_sub_accounts" db:"max_sub_accounts"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// Beneficiary is a saved payee/recipient for payments.
type Beneficiary struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	Name             string     `json:"name" db:"name"`
	IBAN             string     `json:"iban" db:"iban"`
	BIC              string     `json:"bic,omitempty" db:"bic"`
	BankName         string     `json:"bank_name,omitempty" db:"bank_name"`
	DefaultReference string     `json:"default_reference,omitempty" db:"default_reference"`
	IsInternal       bool       `json:"is_internal" db:"is_internal"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// -- API Request/Response Types --

// SubAccountBalance is the balance info returned in account list responses.
type SubAccountBalance struct {
	Available string `json:"available"`
	Pending   string `json:"pending"`
	Total     string `json:"total"`
}

// SubAccountResponse is the API response for a sub-account.
type SubAccountResponse struct {
	ID        string             `json:"id"`
	Currency  string             `json:"currency"`
	IBAN      *string            `json:"iban"`
	BIC       *string            `json:"bic,omitempty"`
	Balance   SubAccountBalance  `json:"balance"`
	IsDefault bool               `json:"is_default"`
}

// AccountResponse is the API response for a user account with sub-accounts.
type AccountResponse struct {
	ID              string               `json:"id"`
	AccountNumber   string               `json:"account_number"`
	Status          string               `json:"status"`
	SubAccounts     []SubAccountResponse `json:"sub_accounts"`
	TotalBalanceEUR string               `json:"total_balance_eur"`
}

// CreateSubAccountRequest matches POST /api/v1/accounts/{id}/sub-accounts.
type CreateSubAccountRequest struct {
	Currency string `json:"currency" binding:"required,len=3"`
}

// -- IBAN Generation --

// BankCode is the TeslaPay bank code used in Lithuanian IBANs.
const BankCode = "TESL"

// BIC is the TeslaPay SWIFT/BIC code.
const BIC = "TESLLT21"

// CountryCode is the ISO 3166 country code for Lithuania.
const CountryCode = "LT"

// GenerateIBAN creates a Lithuanian format IBAN for a given account number.
// Format: LT + check_digits(2) + bank_code(4) + account_number(12)
// Check digits are computed per ISO 13616.
func GenerateIBAN(accountNumber string) string {
	// Pad account number to 12 digits.
	padded := fmt.Sprintf("%012s", accountNumber)

	// The BBAN is bank_code + account_number.
	bban := BankCode + padded

	// Compute check digits per ISO 13616 (mod 97).
	// Rearrange: BBAN + country_code_numeric + "00"
	// LT -> L=21, T=29, so "2129"
	rearranged := bban + "2129" + "00"
	numericStr := convertToNumeric(rearranged)

	n := new(big.Int)
	n.SetString(numericStr, 10)

	mod := new(big.Int)
	mod.Mod(n, big.NewInt(97))

	checkDigits := 98 - mod.Int64()

	return fmt.Sprintf("%s%02d%s", CountryCode, checkDigits, bban)
}

// convertToNumeric replaces letters with their numeric equivalents (A=10, B=11, ..., Z=35).
func convertToNumeric(s string) string {
	var result []byte
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			val := int(c - 'A' + 10)
			result = append(result, []byte(fmt.Sprintf("%d", val))...)
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

// ValidateIBAN performs basic IBAN validation using the mod-97 check.
func ValidateIBAN(iban string) bool {
	if len(iban) < 5 {
		return false
	}

	// Rearrange: move first 4 characters to end.
	rearranged := iban[4:] + iban[:4]
	numericStr := convertToNumeric(rearranged)

	n := new(big.Int)
	n.SetString(numericStr, 10)

	mod := new(big.Int)
	mod.Mod(n, big.NewInt(97))

	return mod.Int64() == 1
}

// GenerateExternalID creates a human-readable external ID in the format TP-XXXXXXXX.
func GenerateExternalID() string {
	id := uuid.New()
	// Use first 8 hex chars of UUID for readability.
	hex := fmt.Sprintf("%X", id[:4])
	return fmt.Sprintf("TP-%s", hex)
}
