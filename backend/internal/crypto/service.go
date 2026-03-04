package crypto

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
	"github.com/teslapay/backend/pkg/events"
)

// cryptoPrices holds hardcoded token prices (EUR) — Phase 4 placeholder.
var cryptoPrices = map[string]CryptoPrice{
	"FUSE": {Symbol: "FUSE", Name: "Fuse", PriceEUR: decimal.NewFromFloat(0.032), PriceUSD: decimal.NewFromFloat(0.035), Change24h: decimal.NewFromFloat(2.5)},
	"USDC": {Symbol: "USDC", Name: "USD Coin", PriceEUR: decimal.NewFromFloat(0.92), PriceUSD: decimal.NewFromFloat(1.00), Change24h: decimal.NewFromFloat(0.01)},
	"USDT": {Symbol: "USDT", Name: "Tether", PriceEUR: decimal.NewFromFloat(0.92), PriceUSD: decimal.NewFromFloat(1.00), Change24h: decimal.NewFromFloat(-0.01)},
	"WETH": {Symbol: "WETH", Name: "Wrapped Ether", PriceEUR: decimal.NewFromFloat(3450.00), PriceUSD: decimal.NewFromFloat(3725.00), Change24h: decimal.NewFromFloat(-1.8)},
}

// defaultTokens defines the tokens initialized for every new wallet.
var defaultTokens = []struct {
	Symbol  string
	Name    string
	Address *string
}{
	{Symbol: "FUSE", Name: "Fuse", Address: nil},
	{Symbol: "USDC", Name: "USD Coin", Address: strPtr("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")},
	{Symbol: "USDT", Name: "Tether", Address: strPtr("0xdac17f958d2ee523a2206206994597c13d831ec7")},
	{Symbol: "WETH", Name: "Wrapped Ether", Address: strPtr("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")},
}

// quoteStore is an in-memory store for active quotes (Phase 4 placeholder).
var quoteStore sync.Map // map[quoteID]*CryptoQuote

// Service implements the crypto business logic.
type Service struct {
	repo     *Repository
	producer *events.Producer
	logger   *zap.Logger
}

// NewService creates a new crypto service.
func NewService(repo *Repository, producer *events.Producer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

// GetOrCreateWallet returns the user's wallet, creating it if it doesn't exist.
// On creation, generates a random hex address and initializes all token balances at 0.
func (s *Service) GetOrCreateWallet(ctx context.Context, userIDStr string) (*CryptoWallet, []*CryptoBalance, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, nil, common.NewValidationError("Invalid user_id", "")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get wallet: %w", err)
	}

	if wallet == nil {
		// Create a new wallet.
		address, err := generateWalletAddress()
		if err != nil {
			return nil, nil, fmt.Errorf("generate wallet address: %w", err)
		}

		wallet = &CryptoWallet{
			ID:        uuid.New(),
			UserID:    userID,
			Address:   address,
			Network:   NetworkFuse,
			Status:    WalletStatusActive,
			CreatedAt: time.Now().UTC(),
		}

		if err := s.repo.CreateWallet(ctx, wallet); err != nil {
			return nil, nil, fmt.Errorf("create wallet: %w", err)
		}

		// Initialize token balances.
		now := time.Now().UTC()
		for _, t := range defaultTokens {
			b := &CryptoBalance{
				WalletID:     wallet.ID,
				TokenSymbol:  t.Symbol,
				TokenName:    t.Name,
				TokenAddress: t.Address,
				Balance:      decimal.Zero,
				UpdatedAt:    now,
			}
			if err := s.repo.CreateBalance(ctx, b); err != nil {
				return nil, nil, fmt.Errorf("create balance for %s: %w", t.Symbol, err)
			}
		}

		if s.logger != nil {
			s.logger.Info("crypto wallet created",
				zap.String("wallet_id", wallet.ID.String()),
				zap.String("user_id", userID.String()),
				zap.String("address", address),
			)
		}

		s.publishCryptoEvent(ctx, wallet.ID.String(), userID.String(), "wallet.created", map[string]any{
			"wallet_id": wallet.ID.String(),
			"user_id":   userID.String(),
			"address":   address,
			"network":   wallet.Network,
		})
	}

	balances, err := s.repo.GetBalances(ctx, wallet.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("get balances: %w", err)
	}

	return wallet, balances, nil
}

// GetBalances returns wallet balances with EUR values calculated from hardcoded prices.
func (s *Service) GetBalances(ctx context.Context, userIDStr string) ([]*CryptoBalance, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, common.NewNotFoundError("Crypto wallet not found", "")
	}

	balances, err := s.repo.GetBalances(ctx, wallet.ID)
	if err != nil {
		return nil, fmt.Errorf("get balances: %w", err)
	}

	return balances, nil
}

// GetPrices returns current token prices.
func (s *Service) GetPrices(ctx context.Context) ([]*CryptoPrice, error) {
	now := time.Now().UTC()
	prices := make([]*CryptoPrice, 0, len(cryptoPrices))
	for _, p := range cryptoPrices {
		cp := p
		cp.UpdatedAt = now
		prices = append(prices, &cp)
	}
	return prices, nil
}

// GetQuote generates a buy/sell quote with 1.5% fee, valid for QuoteTTL seconds.
// For BUY: user pays fiatAmount EUR, receives crypto. cryptoAmount = (fiatAmount - fee) / rate
// For SELL: user sells cryptoAmount tokens, receives fiat. fiatAmount = cryptoAmount * rate - fee
func (s *Service) GetQuote(ctx context.Context, userIDStr, action, symbol, amountStr string) (*CryptoQuote, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	price, ok := cryptoPrices[symbol]
	if !ok {
		return nil, common.NewValidationError("Unsupported token symbol: "+symbol, "")
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, common.NewValidationError("Invalid amount: must be a positive decimal number", "")
	}

	feePct := decimal.NewFromFloat(FeePct)
	rate := price.PriceEUR

	var fiatAmount, cryptoAmount, feeAmount decimal.Decimal

	switch action {
	case ActionBuy:
		// User provides fiatAmount EUR; we deduct the fee and compute crypto received.
		fiatAmount = amount
		feeAmount = fiatAmount.Mul(feePct)
		cryptoAmount = fiatAmount.Sub(feeAmount).Div(rate)
	case ActionSell:
		// User provides cryptoAmount tokens; compute gross fiat then deduct fee.
		cryptoAmount = amount
		grossFiat := cryptoAmount.Mul(rate)
		feeAmount = grossFiat.Mul(feePct)
		fiatAmount = grossFiat.Sub(feeAmount)
	default:
		return nil, common.NewValidationError("action must be buy or sell", "")
	}

	quote := &CryptoQuote{
		ID:           uuid.New().String(),
		UserID:       userID,
		Action:       action,
		TokenSymbol:  symbol,
		FiatAmount:   fiatAmount,
		FiatCurrency: "EUR",
		CryptoAmount: cryptoAmount,
		Rate:         rate,
		FeeAmount:    feeAmount,
		FeePct:       feePct,
		ExpiresAt:    time.Now().UTC().Add(QuoteTTL),
	}

	quoteStore.Store(quote.ID, quote)

	return quote, nil
}

// BuyCrypto executes a buy: credits crypto to the wallet and records the transaction.
func (s *Service) BuyCrypto(ctx context.Context, userIDStr, quoteID string) (*CryptoTransaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	quote, err := s.getValidQuote(quoteID)
	if err != nil {
		return nil, err
	}
	if quote.UserID != userID {
		return nil, common.NewBusinessError("CRYPTO_002", "Quote Expired", "Quote has expired or does not exist")
	}
	if quote.Action != ActionBuy {
		return nil, common.NewValidationError("Quote action must be 'buy'", "")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, common.NewNotFoundError("Crypto wallet not found", "")
	}

	// Credit crypto to wallet.
	if err := s.repo.UpdateBalance(ctx, wallet.ID, quote.TokenSymbol, quote.CryptoAmount); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	// Remove used quote.
	quoteStore.Delete(quoteID)

	now := time.Now().UTC()
	fiatCurrency := quote.FiatCurrency
	fiatAmount := quote.FiatAmount
	rate := quote.Rate
	feeAmount := quote.FeeAmount

	tx := &CryptoTransaction{
		ID:           uuid.New(),
		WalletID:     wallet.ID,
		Type:         TxTypeBuy,
		TokenSymbol:  quote.TokenSymbol,
		Amount:       quote.CryptoAmount,
		FiatAmount:   &fiatAmount,
		FiatCurrency: &fiatCurrency,
		Rate:         &rate,
		FeeAmount:    &feeAmount,
		Status:       TxStatusCompleted,
		CreatedAt:    now,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("crypto buy executed",
			zap.String("wallet_id", wallet.ID.String()),
			zap.String("user_id", userIDStr),
			zap.String("symbol", quote.TokenSymbol),
			zap.String("crypto_amount", quote.CryptoAmount.String()),
			zap.String("fiat_amount", quote.FiatAmount.StringFixed(4)),
		)
	}

	s.publishCryptoEvent(ctx, tx.ID.String(), userIDStr, "crypto.bought", map[string]any{
		"transaction_id": tx.ID.String(),
		"wallet_id":      wallet.ID.String(),
		"user_id":        userIDStr,
		"symbol":         tx.TokenSymbol,
		"amount":         tx.Amount.String(),
		"fiat_amount":    fiatAmount.StringFixed(4),
		"fiat_currency":  fiatCurrency,
	})

	return tx, nil
}

// SellCrypto executes a sell: validates balance, debits crypto, credits fiat (placeholder).
func (s *Service) SellCrypto(ctx context.Context, userIDStr, quoteID string) (*CryptoTransaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	quote, err := s.getValidQuote(quoteID)
	if err != nil {
		return nil, err
	}
	if quote.UserID != userID {
		return nil, common.NewBusinessError("CRYPTO_002", "Quote Expired", "Quote has expired or does not exist")
	}
	if quote.Action != ActionSell {
		return nil, common.NewValidationError("Quote action must be 'sell'", "")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, common.NewNotFoundError("Crypto wallet not found", "")
	}

	// Validate sufficient balance.
	balances, err := s.repo.GetBalances(ctx, wallet.ID)
	if err != nil {
		return nil, fmt.Errorf("get balances: %w", err)
	}
	var currentBalance decimal.Decimal
	for _, b := range balances {
		if b.TokenSymbol == quote.TokenSymbol {
			currentBalance = b.Balance
			break
		}
	}
	if currentBalance.LessThan(quote.CryptoAmount) {
		return nil, common.NewBusinessError(
			"CRYPTO_001",
			"Insufficient Balance",
			fmt.Sprintf("Insufficient %s balance: have %s, need %s", quote.TokenSymbol, currentBalance.String(), quote.CryptoAmount.String()),
		)
	}

	// Debit crypto from wallet (negative delta).
	if err := s.repo.UpdateBalance(ctx, wallet.ID, quote.TokenSymbol, quote.CryptoAmount.Neg()); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	// Remove used quote.
	quoteStore.Delete(quoteID)

	now := time.Now().UTC()
	fiatCurrency := quote.FiatCurrency
	fiatAmount := quote.FiatAmount
	rate := quote.Rate
	feeAmount := quote.FeeAmount

	tx := &CryptoTransaction{
		ID:           uuid.New(),
		WalletID:     wallet.ID,
		Type:         TxTypeSell,
		TokenSymbol:  quote.TokenSymbol,
		Amount:       quote.CryptoAmount,
		FiatAmount:   &fiatAmount,
		FiatCurrency: &fiatCurrency,
		Rate:         &rate,
		FeeAmount:    &feeAmount,
		Status:       TxStatusCompleted,
		CreatedAt:    now,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("crypto sell executed",
			zap.String("wallet_id", wallet.ID.String()),
			zap.String("user_id", userIDStr),
			zap.String("symbol", quote.TokenSymbol),
			zap.String("crypto_amount", quote.CryptoAmount.String()),
			zap.String("fiat_amount", quote.FiatAmount.StringFixed(4)),
		)
	}

	s.publishCryptoEvent(ctx, tx.ID.String(), userIDStr, "crypto.sold", map[string]any{
		"transaction_id": tx.ID.String(),
		"wallet_id":      wallet.ID.String(),
		"user_id":        userIDStr,
		"symbol":         tx.TokenSymbol,
		"amount":         tx.Amount.String(),
		"fiat_amount":    fiatAmount.StringFixed(4),
		"fiat_currency":  fiatCurrency,
	})

	return tx, nil
}

// SendCrypto sends crypto to an external address (placeholder blockchain call).
// Validates sufficient balance, debits wallet, creates pending tx, publishes event.
func (s *Service) SendCrypto(ctx context.Context, userIDStr, symbol, amountStr, recipientAddress string) (*CryptoTransaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	if _, ok := cryptoPrices[symbol]; !ok {
		return nil, common.NewValidationError("Unsupported token symbol: "+symbol, "")
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, common.NewValidationError("Invalid amount: must be a positive decimal number", "")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, common.NewNotFoundError("Crypto wallet not found", "")
	}

	// Validate sufficient balance.
	balances, err := s.repo.GetBalances(ctx, wallet.ID)
	if err != nil {
		return nil, fmt.Errorf("get balances: %w", err)
	}
	var currentBalance decimal.Decimal
	for _, b := range balances {
		if b.TokenSymbol == symbol {
			currentBalance = b.Balance
			break
		}
	}
	if currentBalance.LessThan(amount) {
		return nil, common.NewBusinessError(
			"CRYPTO_001",
			"Insufficient Balance",
			fmt.Sprintf("Insufficient %s balance: have %s, need %s", symbol, currentBalance.String(), amount.String()),
		)
	}

	// Debit wallet.
	if err := s.repo.UpdateBalance(ctx, wallet.ID, symbol, amount.Neg()); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	addr := recipientAddress
	now := time.Now().UTC()

	tx := &CryptoTransaction{
		ID:               uuid.New(),
		WalletID:         wallet.ID,
		Type:             TxTypeSend,
		TokenSymbol:      symbol,
		Amount:           amount,
		Status:           TxStatusPending,
		RecipientAddress: &addr,
		CreatedAt:        now,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("crypto send initiated",
			zap.String("wallet_id", wallet.ID.String()),
			zap.String("user_id", userIDStr),
			zap.String("symbol", symbol),
			zap.String("amount", amount.String()),
			zap.String("recipient", recipientAddress),
		)
	}

	s.publishCryptoEvent(ctx, tx.ID.String(), userIDStr, "crypto.sent", map[string]any{
		"transaction_id":    tx.ID.String(),
		"wallet_id":         wallet.ID.String(),
		"user_id":           userIDStr,
		"symbol":            symbol,
		"amount":            amount.String(),
		"recipient_address": recipientAddress,
	})

	return tx, nil
}

// GetTransactions returns paginated transaction history for a user's wallet.
func (s *Service) GetTransactions(ctx context.Context, userIDStr string, limit, offset int) ([]*CryptoTransaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, common.NewValidationError("Invalid user_id", "")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, common.NewNotFoundError("Crypto wallet not found", "")
	}

	return s.repo.GetTransactions(ctx, wallet.ID, limit, offset)
}

// getValidQuote retrieves a quote from the store and validates it has not expired.
func (s *Service) getValidQuote(quoteID string) (*CryptoQuote, error) {
	val, ok := quoteStore.Load(quoteID)
	if !ok {
		return nil, common.NewBusinessError("CRYPTO_002", "Quote Expired", "Quote has expired or does not exist")
	}
	quote := val.(*CryptoQuote)
	if time.Now().UTC().After(quote.ExpiresAt) {
		quoteStore.Delete(quoteID)
		return nil, common.NewBusinessError("CRYPTO_002", "Quote Expired", "Quote has expired or does not exist")
	}
	return quote, nil
}

// publishCryptoEvent publishes a crypto domain event to Kafka (best-effort).
func (s *Service) publishCryptoEvent(ctx context.Context, aggregateID, actorID, eventType string, data map[string]any) {
	if s.producer == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to marshal crypto event data", zap.Error(err))
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
		Source:        "crypto-service",
		AggregateID:   aggregateID,
		AggregateType: "crypto",
		Data:          payload,
		Metadata: events.EventMetadata{
			CorrelationID: correlationID,
			ActorID:       actorID,
			ActorType:     "user",
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := s.producer.Publish(ctx, events.TopicCryptoEvents, event); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to publish crypto event",
				zap.String("aggregate_id", aggregateID),
				zap.String("event_type", eventType),
				zap.Error(err),
			)
		}
	}
}

// generateWalletAddress generates a random Ethereum-style wallet address (0x + 40 hex chars).
func generateWalletAddress() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return "0x" + hex.EncodeToString(b), nil
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}
