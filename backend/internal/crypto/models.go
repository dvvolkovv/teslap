// Package crypto implements the crypto wallet domain for TeslaPay,
// supporting Fuse Network, USDC, USDT, and WETH token management.
package crypto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Network constants.
const (
	NetworkFuse     = "fuse"
	NetworkEthereum = "ethereum"
)

// Wallet status constants.
const (
	WalletStatusActive = "active"
	WalletStatusFrozen = "frozen"
)

// Transaction type constants.
const (
	TxTypeBuy     = "buy"
	TxTypeSell    = "sell"
	TxTypeSend    = "send"
	TxTypeReceive = "receive"
)

// Transaction status constants.
const (
	TxStatusPending   = "pending"
	TxStatusCompleted = "completed"
	TxStatusFailed    = "failed"
)

// Trade action constants.
const (
	ActionBuy  = "buy"
	ActionSell = "sell"
)

// FeePct is the buy/sell fee percentage (1.5%).
const FeePct = 0.015

// QuoteTTL is how long a quote is valid.
const QuoteTTL = 60 * time.Second

// CryptoWallet represents a user's on-chain wallet.
type CryptoWallet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Address   string    `json:"address" db:"address"`
	Network   string    `json:"network" db:"network"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CryptoBalance represents the token balance for a wallet.
type CryptoBalance struct {
	WalletID     uuid.UUID       `json:"wallet_id" db:"wallet_id"`
	TokenSymbol  string          `json:"token_symbol" db:"token_symbol"`
	TokenName    string          `json:"token_name" db:"token_name"`
	TokenAddress *string         `json:"token_address,omitempty" db:"token_address"`
	Balance      decimal.Decimal `json:"balance" db:"balance"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// CryptoPrice represents the current market price of a token.
type CryptoPrice struct {
	Symbol    string          `json:"symbol"`
	Name      string          `json:"name"`
	PriceEUR  decimal.Decimal `json:"price_eur"`
	PriceUSD  decimal.Decimal `json:"price_usd"`
	Change24h decimal.Decimal `json:"change_24h"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// CryptoQuote is a time-limited buy/sell quote with fee details.
type CryptoQuote struct {
	ID           string          `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	Action       string          `json:"action"`
	TokenSymbol  string          `json:"token_symbol"`
	FiatAmount   decimal.Decimal `json:"fiat_amount"`
	FiatCurrency string          `json:"fiat_currency"`
	CryptoAmount decimal.Decimal `json:"crypto_amount"`
	Rate         decimal.Decimal `json:"rate"`
	FeeAmount    decimal.Decimal `json:"fee_amount"`
	FeePct       decimal.Decimal `json:"fee_pct"`
	ExpiresAt    time.Time       `json:"expires_at"`
}

// CryptoTransaction records a buy, sell, send, or receive operation.
type CryptoTransaction struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	WalletID         uuid.UUID        `json:"wallet_id" db:"wallet_id"`
	TxHash           *string          `json:"tx_hash,omitempty" db:"tx_hash"`
	Type             string           `json:"type" db:"type"`
	TokenSymbol      string           `json:"token_symbol" db:"token_symbol"`
	Amount           decimal.Decimal  `json:"amount" db:"amount"`
	FiatAmount       *decimal.Decimal `json:"fiat_amount,omitempty" db:"fiat_amount"`
	FiatCurrency     *string          `json:"fiat_currency,omitempty" db:"fiat_currency"`
	Rate             *decimal.Decimal `json:"rate,omitempty" db:"rate"`
	FeeAmount        *decimal.Decimal `json:"fee_amount,omitempty" db:"fee_amount"`
	Status           string           `json:"status" db:"status"`
	RecipientAddress *string          `json:"recipient_address,omitempty" db:"recipient_address"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
}

// -- Request types --

// GetQuoteRequest is the query for GET /crypto/quote.
type GetQuoteRequest struct {
	Action string `form:"action" binding:"required,oneof=buy sell"`
	Symbol string `form:"symbol" binding:"required"`
	Amount string `form:"amount" binding:"required"`
}

// BuySellRequest is the request body for POST /crypto/buy and /crypto/sell.
type BuySellRequest struct {
	QuoteID string `json:"quote_id" binding:"required"`
}

// SendCryptoRequest is the request body for POST /crypto/send.
type SendCryptoRequest struct {
	Symbol           string `json:"symbol" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
	RecipientAddress string `json:"recipient_address" binding:"required"`
}

// -- Response types --

// WalletResponse is the API response for wallet operations.
type WalletResponse struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Address   string             `json:"address"`
	Network   string             `json:"network"`
	Status    string             `json:"status"`
	Balances  []*BalanceResponse `json:"balances"`
	CreatedAt string             `json:"created_at"`
}

// BalanceResponse is a single token balance in the wallet response.
type BalanceResponse struct {
	TokenSymbol string `json:"token_symbol"`
	TokenName   string `json:"token_name"`
	Balance     string `json:"balance"`
	ValueEUR    string `json:"value_eur"`
}

// TransactionResponse is the API response for a single crypto transaction.
type TransactionResponse struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	TokenSymbol      string `json:"token_symbol"`
	Amount           string `json:"amount"`
	FiatAmount       string `json:"fiat_amount,omitempty"`
	FiatCurrency     string `json:"fiat_currency,omitempty"`
	Status           string `json:"status"`
	RecipientAddress string `json:"recipient_address,omitempty"`
	CreatedAt        string `json:"created_at"`
}

// ListTransactionsResponse wraps a paginated list of transactions.
type ListTransactionsResponse struct {
	Data  []*TransactionResponse `json:"data"`
	Total int                    `json:"total"`
}

// toTransactionResponse converts a CryptoTransaction to its API response form.
func toTransactionResponse(tx *CryptoTransaction) *TransactionResponse {
	resp := &TransactionResponse{
		ID:          tx.ID.String(),
		Type:        tx.Type,
		TokenSymbol: tx.TokenSymbol,
		Amount:      tx.Amount.String(),
		Status:      tx.Status,
		CreatedAt:   tx.CreatedAt.Format(time.RFC3339),
	}
	if tx.FiatAmount != nil {
		resp.FiatAmount = tx.FiatAmount.StringFixed(4)
	}
	if tx.FiatCurrency != nil {
		resp.FiatCurrency = *tx.FiatCurrency
	}
	if tx.RecipientAddress != nil {
		resp.RecipientAddress = *tx.RecipientAddress
	}
	return resp
}

// toWalletResponse converts a CryptoWallet and its balances to the API response form.
// prices is the current price map used to compute EUR values.
func toWalletResponse(w *CryptoWallet, balances []*CryptoBalance, prices map[string]CryptoPrice) *WalletResponse {
	resp := &WalletResponse{
		ID:        w.ID.String(),
		UserID:    w.UserID.String(),
		Address:   w.Address,
		Network:   w.Network,
		Status:    w.Status,
		Balances:  make([]*BalanceResponse, 0, len(balances)),
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
	}
	for _, b := range balances {
		br := &BalanceResponse{
			TokenSymbol: b.TokenSymbol,
			TokenName:   b.TokenName,
			Balance:     b.Balance.String(),
			ValueEUR:    "0.00",
		}
		if p, ok := prices[b.TokenSymbol]; ok {
			valueEUR := b.Balance.Mul(p.PriceEUR)
			br.ValueEUR = valueEUR.StringFixed(2)
		}
		resp.Balances = append(resp.Balances, br)
	}
	return resp
}
