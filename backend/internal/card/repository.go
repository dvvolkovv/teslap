package card

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the card service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new card repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreateCard inserts a new card record.
func (r *Repository) CreateCard(ctx context.Context, c *Card) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO cards
			(id, account_id, sub_account_id, card_number_encrypted, last_four,
			 expiry_month, expiry_year, cvv_hash, cardholder_name, type, status,
			 daily_limit, monthly_limit, daily_spent, monthly_spent,
			 is_contactless, is_online, is_atm, allowed_countries, blocked_mcc_codes,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		        $16, $17, $18, $19, $20, $21, $22)
	`,
		c.ID, c.AccountID, c.SubAccountID, c.CardNumberEncrypted, c.LastFour,
		c.ExpiryMonth, c.ExpiryYear, c.CVVHash, c.CardholderName, c.Type, c.Status,
		c.DailyLimit, c.MonthlyLimit, c.DailySpent, c.MonthlySpent,
		c.IsContactless, c.IsOnline, c.IsATM, c.AllowedCountries, c.BlockedMCCCodes,
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card: %w", err)
	}
	return nil
}

// GetCardByID fetches a single card by its ID.
func (r *Repository) GetCardByID(ctx context.Context, id uuid.UUID) (*Card, error) {
	var c Card
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, account_id, sub_account_id, card_number_encrypted, last_four,
		       expiry_month, expiry_year, cvv_hash, cardholder_name, type, status,
		       daily_limit, monthly_limit, daily_spent, monthly_spent,
		       is_contactless, is_online, is_atm, allowed_countries, blocked_mcc_codes,
		       created_at, updated_at
		FROM cards
		WHERE id = $1
	`, id).Scan(
		&c.ID, &c.AccountID, &c.SubAccountID, &c.CardNumberEncrypted, &c.LastFour,
		&c.ExpiryMonth, &c.ExpiryYear, &c.CVVHash, &c.CardholderName, &c.Type, &c.Status,
		&c.DailyLimit, &c.MonthlyLimit, &c.DailySpent, &c.MonthlySpent,
		&c.IsContactless, &c.IsOnline, &c.IsATM, &c.AllowedCountries, &c.BlockedMCCCodes,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query card by id: %w", err)
	}
	return &c, nil
}

// GetCardsByAccountID returns all cards for an account.
func (r *Repository) GetCardsByAccountID(ctx context.Context, accountID uuid.UUID) ([]*Card, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, account_id, sub_account_id, card_number_encrypted, last_four,
		       expiry_month, expiry_year, cvv_hash, cardholder_name, type, status,
		       daily_limit, monthly_limit, daily_spent, monthly_spent,
		       is_contactless, is_online, is_atm, allowed_countries, blocked_mcc_codes,
		       created_at, updated_at
		FROM cards
		WHERE account_id = $1
		ORDER BY created_at DESC
	`, accountID)
	if err != nil {
		return nil, fmt.Errorf("query cards by account: %w", err)
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		var c Card
		if err := rows.Scan(
			&c.ID, &c.AccountID, &c.SubAccountID, &c.CardNumberEncrypted, &c.LastFour,
			&c.ExpiryMonth, &c.ExpiryYear, &c.CVVHash, &c.CardholderName, &c.Type, &c.Status,
			&c.DailyLimit, &c.MonthlyLimit, &c.DailySpent, &c.MonthlySpent,
			&c.IsContactless, &c.IsOnline, &c.IsATM, &c.AllowedCountries, &c.BlockedMCCCodes,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan card: %w", err)
		}
		cards = append(cards, &c)
	}
	return cards, rows.Err()
}

// UpdateCardStatus updates a card's status and updated_at timestamp.
func (r *Repository) UpdateCardStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE cards SET status = $2, updated_at = NOW() WHERE id = $1
	`, id, status)
	if err != nil {
		return fmt.Errorf("update card status: %w", err)
	}
	return nil
}

// UpdateCardControls updates spending limits and permission flags.
func (r *Repository) UpdateCardControls(ctx context.Context, id uuid.UUID, update *CardControlsUpdate) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE cards SET
			daily_limit = $2,
			monthly_limit = $3,
			is_contactless = $4,
			is_online = $5,
			is_atm = $6,
			allowed_countries = $7,
			blocked_mcc_codes = $8,
			updated_at = NOW()
		WHERE id = $1
	`, id, update.DailyLimit, update.MonthlyLimit, update.IsContactless, update.IsOnline,
		update.IsATM, update.AllowedCountries, update.BlockedMCCCodes)
	if err != nil {
		return fmt.Errorf("update card controls: %w", err)
	}
	return nil
}

// UpdateCardSpending updates the daily and monthly spent amounts.
func (r *Repository) UpdateCardSpending(ctx context.Context, id uuid.UUID, dailySpent, monthlySpent decimal.Decimal) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE cards SET daily_spent = $2, monthly_spent = $3, updated_at = NOW() WHERE id = $1
	`, id, dailySpent, monthlySpent)
	if err != nil {
		return fmt.Errorf("update card spending: %w", err)
	}
	return nil
}

// CreateCardTransaction inserts a new card transaction record.
func (r *Repository) CreateCardTransaction(ctx context.Context, tx *CardTransaction) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO card_transactions
			(id, card_id, merchant_name, merchant_mcc, amount, currency,
			 original_amount, original_currency, status, type, country,
			 authorization_code, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`,
		tx.ID, tx.CardID, tx.MerchantName, tx.MerchantMCC, tx.Amount, tx.Currency,
		tx.OriginalAmount, tx.OriginalCurrency, tx.Status, tx.Type, tx.Country,
		tx.AuthorizationCode, tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert card transaction: %w", err)
	}
	return nil
}

// GetCardTransactions returns paginated card transactions for a card.
func (r *Repository) GetCardTransactions(ctx context.Context, cardID uuid.UUID, limit, offset int) ([]*CardTransaction, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, card_id, merchant_name, merchant_mcc, amount, currency,
		       original_amount, original_currency, status, type, country,
		       authorization_code, created_at
		FROM card_transactions
		WHERE card_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, cardID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query card transactions: %w", err)
	}
	defer rows.Close()

	var txs []*CardTransaction
	for rows.Next() {
		var tx CardTransaction
		if err := rows.Scan(
			&tx.ID, &tx.CardID, &tx.MerchantName, &tx.MerchantMCC, &tx.Amount, &tx.Currency,
			&tx.OriginalAmount, &tx.OriginalCurrency, &tx.Status, &tx.Type, &tx.Country,
			&tx.AuthorizationCode, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan card transaction: %w", err)
		}
		txs = append(txs, &tx)
	}
	return txs, rows.Err()
}
