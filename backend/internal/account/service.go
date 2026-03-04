package account

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// SupportedCurrencies defines which currencies can be used for sub-accounts.
var SupportedCurrencies = map[string]bool{
	"EUR": true,
	"USD": true,
	"GBP": true,
	"PLN": true,
	"CHF": true,
}

// Service implements account management business logic.
type Service struct {
	repo   *Repository
	logger *zap.Logger
}

// NewService creates a new account service.
func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// CreateAccount creates a new account for a user with a default EUR sub-account.
// Called after successful registration and KYC approval.
func (s *Service) CreateAccount(ctx context.Context, userID uuid.UUID) (*Account, *SubAccount, error) {
	// Generate a unique account number.
	accountNumber := fmt.Sprintf("TP%010d", time.Now().UnixNano()%10000000000)

	account := &Account{
		ID:            uuid.New(),
		UserID:        userID,
		AccountNumber: accountNumber,
		Status:        "active",
		OpenedAt:      time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := s.repo.CreateAccount(ctx, account); err != nil {
		return nil, nil, fmt.Errorf("create account: %w", err)
	}

	// Create a default EUR sub-account with an IBAN.
	subAccount, err := s.createSubAccount(ctx, account.ID, "EUR", true)
	if err != nil {
		return nil, nil, fmt.Errorf("create default sub-account: %w", err)
	}

	s.logger.Info("account created",
		zap.String("user_id", userID.String()),
		zap.String("account_id", account.ID.String()),
		zap.String("iban", safeDeref(subAccount.IBAN)),
	)

	return account, subAccount, nil
}

// GetAccount retrieves an account with its sub-accounts by account ID.
func (s *Service) GetAccount(ctx context.Context, accountID uuid.UUID) (*Account, []SubAccount, error) {
	account, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, nil, fmt.Errorf("get account: %w", err)
	}
	if account == nil {
		return nil, nil, common.NewNotFoundError("Account not found", "")
	}

	subAccounts, err := s.repo.GetSubAccountsByAccountID(ctx, accountID)
	if err != nil {
		return nil, nil, fmt.Errorf("get sub-accounts: %w", err)
	}

	return account, subAccounts, nil
}

// GetAccountsByUserID retrieves all accounts for a user.
func (s *Service) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	return s.repo.GetAccountsByUserID(ctx, userID)
}

// OpenSubAccount creates a new sub-account in the specified currency.
func (s *Service) OpenSubAccount(ctx context.Context, accountID uuid.UUID, currency string) (*SubAccount, error) {
	if !SupportedCurrencies[currency] {
		return nil, common.NewValidationError(
			fmt.Sprintf("Currency %s is not supported. Supported: EUR, USD, GBP, PLN, CHF", currency),
			"",
		)
	}

	// Check that the account exists and is active.
	account, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("lookup account: %w", err)
	}
	if account == nil {
		return nil, common.NewNotFoundError("Account not found", "")
	}
	if account.Status != "active" {
		return nil, common.NewBusinessError("ACCOUNT_001", "Account Inactive",
			"Cannot create sub-account on an inactive account")
	}

	// Check if a sub-account already exists for this currency.
	existing, err := s.repo.GetSubAccountByAccountAndCurrency(ctx, accountID, currency)
	if err != nil {
		return nil, fmt.Errorf("check existing sub-account: %w", err)
	}
	if existing != nil {
		return nil, common.NewConflictError(
			fmt.Sprintf("A %s sub-account already exists for this account", currency),
		)
	}

	// Check tier limit on sub-accounts.
	subAccounts, err := s.repo.GetSubAccountsByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("count sub-accounts: %w", err)
	}
	// Default max is 5; in production, check against user's tier.
	if len(subAccounts) >= 5 {
		return nil, common.NewBusinessError("ACCOUNT_002", "Sub-Account Limit Reached",
			"Maximum number of sub-accounts reached for your tier")
	}

	subAccount, err := s.createSubAccount(ctx, accountID, currency, false)
	if err != nil {
		return nil, fmt.Errorf("create sub-account: %w", err)
	}

	s.logger.Info("sub-account created",
		zap.String("account_id", accountID.String()),
		zap.String("currency", currency),
		zap.String("sub_account_id", subAccount.ID.String()),
	)

	return subAccount, nil
}

// UpdateTier updates a user's account tier.
func (s *Service) UpdateTier(ctx context.Context, userID uuid.UUID, tierName string) error {
	tier, err := s.repo.GetTierByName(ctx, tierName)
	if err != nil {
		return fmt.Errorf("lookup tier: %w", err)
	}
	if tier == nil {
		return common.NewNotFoundError("Tier not found: "+tierName, "")
	}

	if err := s.repo.UpdateUserTier(ctx, userID, tier.ID); err != nil {
		return fmt.Errorf("update tier: %w", err)
	}

	s.logger.Info("user tier updated",
		zap.String("user_id", userID.String()),
		zap.String("tier", tierName),
	)

	return nil
}

// createSubAccount is the internal helper that creates a sub-account.
func (s *Service) createSubAccount(ctx context.Context, accountID uuid.UUID, currency string, isDefault bool) (*SubAccount, error) {
	ledgerAccountID := uuid.New() // In production, created via ledger service gRPC call.

	subAccount := &SubAccount{
		ID:              uuid.New(),
		AccountID:       accountID,
		Currency:        currency,
		LedgerAccountID: ledgerAccountID,
		Status:          "active",
		IsDefault:       isDefault,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	// Generate IBAN for EUR accounts.
	if currency == "EUR" {
		accountNumStr := fmt.Sprintf("%012d", time.Now().UnixNano()%1000000000000)
		iban := GenerateIBAN(accountNumStr)
		bic := BIC
		subAccount.IBAN = &iban
		subAccount.BIC = &bic
	}

	if err := s.repo.CreateSubAccount(ctx, subAccount); err != nil {
		return nil, fmt.Errorf("insert sub-account: %w", err)
	}

	return subAccount, nil
}

func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
