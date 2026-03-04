package crypto

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the crypto service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new crypto repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreateWallet inserts a new crypto wallet.
func (r *Repository) CreateWallet(ctx context.Context, w *CryptoWallet) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO crypto_wallets (id, user_id, address, network, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, w.ID, w.UserID, w.Address, w.Network, w.Status, w.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert crypto wallet: %w", err)
	}
	return nil
}

// GetWalletByUserID finds a wallet by user_id. Returns nil, nil if not found.
func (r *Repository) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*CryptoWallet, error) {
	var w CryptoWallet
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, address, network, status, created_at
		FROM crypto_wallets
		WHERE user_id = $1
	`, userID).Scan(&w.ID, &w.UserID, &w.Address, &w.Network, &w.Status, &w.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query crypto wallet by user_id: %w", err)
	}
	return &w, nil
}

// GetBalances returns all token balances for a wallet.
func (r *Repository) GetBalances(ctx context.Context, walletID uuid.UUID) ([]*CryptoBalance, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT wallet_id, token_symbol, token_name, token_address, balance, updated_at
		FROM crypto_balances
		WHERE wallet_id = $1
		ORDER BY token_symbol
	`, walletID)
	if err != nil {
		return nil, fmt.Errorf("query crypto balances: %w", err)
	}
	defer rows.Close()

	var balances []*CryptoBalance
	for rows.Next() {
		var b CryptoBalance
		if err := rows.Scan(&b.WalletID, &b.TokenSymbol, &b.TokenName, &b.TokenAddress, &b.Balance, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan crypto balance: %w", err)
		}
		balances = append(balances, &b)
	}
	return balances, rows.Err()
}

// CreateBalance inserts a token balance row (for wallet initialization).
func (r *Repository) CreateBalance(ctx context.Context, b *CryptoBalance) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO crypto_balances (wallet_id, token_symbol, token_name, token_address, balance, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, b.WalletID, b.TokenSymbol, b.TokenName, b.TokenAddress, b.Balance, b.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert crypto balance: %w", err)
	}
	return nil
}

// UpdateBalance updates (adds delta) the balance for a specific token.
func (r *Repository) UpdateBalance(ctx context.Context, walletID uuid.UUID, tokenSymbol string, delta decimal.Decimal) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE crypto_balances
		SET balance = balance + $3, updated_at = NOW()
		WHERE wallet_id = $1 AND token_symbol = $2
	`, walletID, tokenSymbol, delta)
	if err != nil {
		return fmt.Errorf("update crypto balance: %w", err)
	}
	return nil
}

// CreateTransaction inserts a new crypto transaction record.
func (r *Repository) CreateTransaction(ctx context.Context, tx *CryptoTransaction) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO crypto_transactions
			(id, wallet_id, tx_hash, type, token_symbol, amount, fiat_amount, fiat_currency, rate, fee_amount, status, recipient_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`,
		tx.ID, tx.WalletID, tx.TxHash, tx.Type, tx.TokenSymbol,
		tx.Amount, tx.FiatAmount, tx.FiatCurrency, tx.Rate, tx.FeeAmount,
		tx.Status, tx.RecipientAddress, tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert crypto transaction: %w", err)
	}
	return nil
}

// GetTransactions returns paginated transactions for a wallet.
func (r *Repository) GetTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*CryptoTransaction, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, wallet_id, tx_hash, type, token_symbol, amount, fiat_amount, fiat_currency, rate, fee_amount, status, recipient_address, created_at
		FROM crypto_transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query crypto transactions: %w", err)
	}
	defer rows.Close()

	var txs []*CryptoTransaction
	for rows.Next() {
		var tx CryptoTransaction
		if err := rows.Scan(
			&tx.ID, &tx.WalletID, &tx.TxHash, &tx.Type, &tx.TokenSymbol,
			&tx.Amount, &tx.FiatAmount, &tx.FiatCurrency, &tx.Rate, &tx.FeeAmount,
			&tx.Status, &tx.RecipientAddress, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan crypto transaction: %w", err)
		}
		txs = append(txs, &tx)
	}
	return txs, rows.Err()
}
