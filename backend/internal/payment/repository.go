package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the payment service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new payment repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreatePayment inserts a new payment record.
func (r *Repository) CreatePayment(ctx context.Context, p *Payment) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO payments
			(id, sender_account_id, recipient_account_id, recipient_iban, recipient_name,
			 amount, currency, type, status, reference, description, idempotency_key,
			 fee_amount, fee_currency, fx_rate, fx_from_currency, fx_to_currency,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`,
		p.ID, p.SenderAccountID, p.RecipientAccountID, p.RecipientIBAN, p.RecipientName,
		p.Amount, p.Currency, p.Type, p.Status, p.Reference, p.Description, p.IdempotencyKey,
		p.FeeAmount, p.FeeCurrency, p.FXRate, p.FXFromCurrency, p.FXToCurrency,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	return nil
}

// GetPaymentByID fetches a single payment by its ID.
func (r *Repository) GetPaymentByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
	var p Payment
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, sender_account_id, recipient_account_id, recipient_iban, recipient_name,
		       amount, currency, type, status, reference, description, idempotency_key,
		       fee_amount, fee_currency, fx_rate, fx_from_currency, fx_to_currency,
		       created_at, updated_at
		FROM payments
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.SenderAccountID, &p.RecipientAccountID, &p.RecipientIBAN, &p.RecipientName,
		&p.Amount, &p.Currency, &p.Type, &p.Status, &p.Reference, &p.Description, &p.IdempotencyKey,
		&p.FeeAmount, &p.FeeCurrency, &p.FXRate, &p.FXFromCurrency, &p.FXToCurrency,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query payment by id: %w", err)
	}
	return &p, nil
}

// GetPaymentsByAccountID returns paginated payments for a sender account.
func (r *Repository) GetPaymentsByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*Payment, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, sender_account_id, recipient_account_id, recipient_iban, recipient_name,
		       amount, currency, type, status, reference, description, idempotency_key,
		       fee_amount, fee_currency, fx_rate, fx_from_currency, fx_to_currency,
		       created_at, updated_at
		FROM payments
		WHERE sender_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query payments by account: %w", err)
	}
	defer rows.Close()

	var payments []*Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(
			&p.ID, &p.SenderAccountID, &p.RecipientAccountID, &p.RecipientIBAN, &p.RecipientName,
			&p.Amount, &p.Currency, &p.Type, &p.Status, &p.Reference, &p.Description, &p.IdempotencyKey,
			&p.FeeAmount, &p.FeeCurrency, &p.FXRate, &p.FXFromCurrency, &p.FXToCurrency,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, &p)
	}
	return payments, rows.Err()
}

// UpdatePaymentStatus updates the status of a payment.
func (r *Repository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE payments SET status = $2, updated_at = NOW() WHERE id = $1
	`, id, status)
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}
	return nil
}

// GetPaymentByIdempotencyKey finds a payment by its idempotency key.
func (r *Repository) GetPaymentByIdempotencyKey(ctx context.Context, key string) (*Payment, error) {
	var p Payment
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, sender_account_id, recipient_account_id, recipient_iban, recipient_name,
		       amount, currency, type, status, reference, description, idempotency_key,
		       fee_amount, fee_currency, fx_rate, fx_from_currency, fx_to_currency,
		       created_at, updated_at
		FROM payments
		WHERE idempotency_key = $1
	`, key).Scan(
		&p.ID, &p.SenderAccountID, &p.RecipientAccountID, &p.RecipientIBAN, &p.RecipientName,
		&p.Amount, &p.Currency, &p.Type, &p.Status, &p.Reference, &p.Description, &p.IdempotencyKey,
		&p.FeeAmount, &p.FeeCurrency, &p.FXRate, &p.FXFromCurrency, &p.FXToCurrency,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query payment by idempotency key: %w", err)
	}
	return &p, nil
}

// CreateScheduledPayment inserts a new scheduled payment.
func (r *Repository) CreateScheduledPayment(ctx context.Context, sp *ScheduledPayment) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO scheduled_payments
			(id, account_id, recipient_account_id, recipient_iban, recipient_name,
			 amount, currency, type, schedule_type, reference, description,
			 is_active, next_execution, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`,
		sp.ID, sp.AccountID, sp.RecipientAccountID, sp.RecipientIBAN, sp.RecipientName,
		sp.Amount, sp.Currency, sp.Type, sp.ScheduleType, sp.Reference, sp.Description,
		sp.IsActive, sp.NextExecution, sp.CreatedAt, sp.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert scheduled payment: %w", err)
	}
	return nil
}

// GetScheduledPayments returns all scheduled payments for an account.
func (r *Repository) GetScheduledPayments(ctx context.Context, accountID uuid.UUID) ([]*ScheduledPayment, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, account_id, recipient_account_id, recipient_iban, recipient_name,
		       amount, currency, type, schedule_type, reference, description,
		       is_active, next_execution, last_execution, created_at, updated_at
		FROM scheduled_payments
		WHERE account_id = $1
		ORDER BY created_at DESC
	`, accountID)
	if err != nil {
		return nil, fmt.Errorf("query scheduled payments: %w", err)
	}
	defer rows.Close()

	var result []*ScheduledPayment
	for rows.Next() {
		var sp ScheduledPayment
		if err := rows.Scan(
			&sp.ID, &sp.AccountID, &sp.RecipientAccountID, &sp.RecipientIBAN, &sp.RecipientName,
			&sp.Amount, &sp.Currency, &sp.Type, &sp.ScheduleType, &sp.Reference, &sp.Description,
			&sp.IsActive, &sp.NextExecution, &sp.LastExecution, &sp.CreatedAt, &sp.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan scheduled payment: %w", err)
		}
		result = append(result, &sp)
	}
	return result, rows.Err()
}
