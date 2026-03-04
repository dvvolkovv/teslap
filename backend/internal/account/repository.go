package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the account service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new account repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreateAccount inserts a new account record.
func (r *Repository) CreateAccount(ctx context.Context, account *Account) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO accounts (id, user_id, account_number, status, opened_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, account.ID, account.UserID, account.AccountNumber,
		account.Status, account.OpenedAt, account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert account: %w", err)
	}
	return nil
}

// GetAccountByID retrieves an account by its ID.
func (r *Repository) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	var a Account
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, account_number, status, opened_at, closed_at, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(
		&a.ID, &a.UserID, &a.AccountNumber, &a.Status,
		&a.OpenedAt, &a.ClosedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query account: %w", err)
	}
	return &a, nil
}

// GetAccountsByUserID retrieves all accounts for a user.
func (r *Repository) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, user_id, account_number, status, opened_at, closed_at, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.AccountNumber, &a.Status,
			&a.OpenedAt, &a.ClosedAt, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

// CreateSubAccount inserts a new sub-account record.
func (r *Repository) CreateSubAccount(ctx context.Context, sa *SubAccount) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO sub_accounts
			(id, account_id, currency, iban, bic, ledger_account_id, status, is_default,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		sa.ID, sa.AccountID, sa.Currency, sa.IBAN, sa.BIC,
		sa.LedgerAccountID, sa.Status, sa.IsDefault,
		sa.CreatedAt, sa.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert sub-account: %w", err)
	}
	return nil
}

// GetSubAccountsByAccountID retrieves all sub-accounts for a given account.
func (r *Repository) GetSubAccountsByAccountID(ctx context.Context, accountID uuid.UUID) ([]SubAccount, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, account_id, currency, iban, bic, ledger_account_id,
		       status, is_default, created_at, updated_at
		FROM sub_accounts
		WHERE account_id = $1
		ORDER BY is_default DESC, created_at
	`, accountID)
	if err != nil {
		return nil, fmt.Errorf("query sub-accounts: %w", err)
	}
	defer rows.Close()

	var subAccounts []SubAccount
	for rows.Next() {
		var sa SubAccount
		if err := rows.Scan(
			&sa.ID, &sa.AccountID, &sa.Currency, &sa.IBAN, &sa.BIC,
			&sa.LedgerAccountID, &sa.Status, &sa.IsDefault,
			&sa.CreatedAt, &sa.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan sub-account: %w", err)
		}
		subAccounts = append(subAccounts, sa)
	}
	return subAccounts, rows.Err()
}

// GetSubAccountByAccountAndCurrency retrieves a sub-account by account ID and currency.
func (r *Repository) GetSubAccountByAccountAndCurrency(
	ctx context.Context, accountID uuid.UUID, currency string,
) (*SubAccount, error) {
	var sa SubAccount
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, account_id, currency, iban, bic, ledger_account_id,
		       status, is_default, created_at, updated_at
		FROM sub_accounts
		WHERE account_id = $1 AND currency = $2
	`, accountID, currency).Scan(
		&sa.ID, &sa.AccountID, &sa.Currency, &sa.IBAN, &sa.BIC,
		&sa.LedgerAccountID, &sa.Status, &sa.IsDefault,
		&sa.CreatedAt, &sa.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query sub-account: %w", err)
	}
	return &sa, nil
}

// GetTierByName retrieves an account tier by its name.
func (r *Repository) GetTierByName(ctx context.Context, name string) (*AccountTier, error) {
	var t AccountTier
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, monthly_fee, daily_transfer_limit, monthly_transfer_limit,
		       daily_card_limit, monthly_atm_limit, free_atm_withdrawals,
		       fx_markup_percent, max_sub_accounts, created_at, updated_at
		FROM account_tiers
		WHERE name = $1
	`, name).Scan(
		&t.ID, &t.Name, &t.MonthlyFee, &t.DailyTransferLimit,
		&t.MonthlyTransferLimit, &t.DailyCardLimit, &t.MonthlyATMLimit,
		&t.FreeATMWithdrawals, &t.FXMarkupPercent, &t.MaxSubAccounts,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query tier: %w", err)
	}
	return &t, nil
}

// UpdateUserTier updates the tier for a given user.
func (r *Repository) UpdateUserTier(ctx context.Context, userID, tierID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE users SET tier_id = $2, updated_at = NOW() WHERE id = $1
	`, userID, tierID)
	return err
}

// GetUserByID retrieves a user profile by ID.
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, external_id, tier_id, first_name, last_name, date_of_birth,
		       nationality, tax_residency, address_line1, address_line2,
		       city, postal_code, country, language, kyc_status, kyc_level,
		       risk_score, status, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&u.ID, &u.ExternalID, &u.TierID, &u.FirstName, &u.LastName,
		&u.DateOfBirth, &u.Nationality, &u.TaxResidency,
		&u.AddressLine1, &u.AddressLine2, &u.City, &u.PostalCode,
		&u.Country, &u.Language, &u.KYCStatus, &u.KYCLevel,
		&u.RiskScore, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	return &u, nil
}
