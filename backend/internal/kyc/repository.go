package kyc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/teslapay/backend/pkg/database"
)

// Repository provides data access for the KYC service.
type Repository struct {
	db     *database.DB
	logger *zap.Logger
}

// NewRepository creates a new KYC repository.
func NewRepository(db *database.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

// CreateKYCRecord inserts a new KYC record.
func (r *Repository) CreateKYCRecord(ctx context.Context, record *KYCRecord) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO kyc_records
			(id, user_id, provider, applicant_id, level, status,
			 review_result, rejection_reasons, documents_provided,
			 verified_at, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`,
		record.ID, record.UserID, record.Provider, record.ApplicantID, record.Level, record.Status,
		record.ReviewResult, record.RejectionReasons, record.DocumentsProvided,
		record.VerifiedAt, record.ExpiresAt, record.CreatedAt, record.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert kyc record: %w", err)
	}
	return nil
}

// GetKYCByUserID fetches the most recent KYC record for a user.
func (r *Repository) GetKYCByUserID(ctx context.Context, userID uuid.UUID) (*KYCRecord, error) {
	var rec KYCRecord
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, provider, applicant_id, level, status,
		       review_result, rejection_reasons, documents_provided,
		       verified_at, expires_at, created_at, updated_at
		FROM kyc_records
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(
		&rec.ID, &rec.UserID, &rec.Provider, &rec.ApplicantID, &rec.Level, &rec.Status,
		&rec.ReviewResult, &rec.RejectionReasons, &rec.DocumentsProvided,
		&rec.VerifiedAt, &rec.ExpiresAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query kyc by user_id: %w", err)
	}
	return &rec, nil
}

// GetKYCByApplicantID fetches a KYC record by its Sumsub applicant ID.
func (r *Repository) GetKYCByApplicantID(ctx context.Context, applicantID string) (*KYCRecord, error) {
	var rec KYCRecord
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, provider, applicant_id, level, status,
		       review_result, rejection_reasons, documents_provided,
		       verified_at, expires_at, created_at, updated_at
		FROM kyc_records
		WHERE applicant_id = $1
	`, applicantID).Scan(
		&rec.ID, &rec.UserID, &rec.Provider, &rec.ApplicantID, &rec.Level, &rec.Status,
		&rec.ReviewResult, &rec.RejectionReasons, &rec.DocumentsProvided,
		&rec.VerifiedAt, &rec.ExpiresAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query kyc by applicant_id: %w", err)
	}
	return &rec, nil
}

// UpdateKYCStatus updates the status and review result of a KYC record.
func (r *Repository) UpdateKYCStatus(ctx context.Context, applicantID string, status string, reviewResult json.RawMessage) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE kyc_records
		SET status = $2, review_result = $3, updated_at = NOW()
		WHERE applicant_id = $1
	`, applicantID, status, reviewResult)
	if err != nil {
		return fmt.Errorf("update kyc status: %w", err)
	}
	return nil
}
