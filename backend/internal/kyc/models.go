package kyc

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// KYCRecord stores KYC verification state for a user.
type KYCRecord struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	UserID            uuid.UUID       `json:"user_id" db:"user_id"`
	Provider          string          `json:"provider" db:"provider"`
	ApplicantID       string          `json:"applicant_id,omitempty" db:"applicant_id"`
	Level             string          `json:"level" db:"level"`
	Status            string          `json:"status" db:"status"`
	ReviewResult      json.RawMessage `json:"review_result,omitempty" db:"review_result"`
	RejectionReasons  []string        `json:"rejection_reasons,omitempty" db:"rejection_reasons"`
	DocumentsProvided []string        `json:"documents_provided,omitempty" db:"documents_provided"`
	VerifiedAt        *time.Time      `json:"verified_at,omitempty" db:"verified_at"`
	ExpiresAt         *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
}

// KYC level constants.
const (
	LevelBasic    = "basic"
	LevelEnhanced = "enhanced"
	LevelFull     = "full"
)

// KYC status constants.
const (
	StatusPending  = "pending"
	StatusInReview = "in_review"
	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusExpired  = "expired"
)

// -- Request / Response types --

// KYCStartRequest matches POST /kyc/start.
type KYCStartRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Level  string `json:"level" binding:"required,oneof=basic enhanced full"`
}

// KYCStartResponse is returned after initiating KYC.
type KYCStartResponse struct {
	ApplicantID string `json:"applicant_id"`
	SDKToken    string `json:"sdk_token"`
	FlowURL     string `json:"flow_url"`
}

// KYCStatusResponse is returned for GET /kyc/status.
type KYCStatusResponse struct {
	Level        string      `json:"level"`
	Status       string      `json:"status"`
	ReviewResult interface{} `json:"review_result,omitempty"`
	VerifiedAt   *time.Time  `json:"verified_at,omitempty"`
}

// KYCWebhookPayload is the incoming payload from Sumsub webhooks.
type KYCWebhookPayload struct {
	ApplicantID     string          `json:"applicantId"`
	ReviewStatus    string          `json:"reviewStatus"`
	ReviewResult    json.RawMessage `json:"reviewResult,omitempty"`
	RejectionLabels []string        `json:"rejectionLabels,omitempty"`
}
