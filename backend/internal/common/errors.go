// Package common provides shared types, error handling, and utilities
// used across all TeslaPay microservices.
package common

import (
	"fmt"
	"net/http"
)

// ProblemDetail implements RFC 7807 Problem Details for HTTP APIs.
// All error responses from TeslaPay services follow this format.
type ProblemDetail struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Instance  string `json:"instance,omitempty"`
	ErrorCode string `json:"error_code"`
	TraceID   string `json:"trace_id,omitempty"`
}

func (p *ProblemDetail) Error() string {
	return fmt.Sprintf("[%s] %s: %s", p.ErrorCode, p.Title, p.Detail)
}

// Standard error codes as defined in the API contracts.
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeAuthMissing      = "AUTH_001"
	ErrCodeAuthExpired       = "AUTH_002"
	ErrCodeAuthForbidden    = "AUTH_003"
	ErrCodeAuthNewDevice    = "AUTH_004"
	ErrCodeSCARequired      = "SCA_001"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInsufficientFunds = "PAY_001"
	ErrCodeTransferLimit    = "PAY_002"
	ErrCodeInvalidIBAN      = "PAY_003"
	ErrCodeKYCRequired      = "KYC_001"
	ErrCodeCardFrozen       = "CARD_001"
	ErrCodeRateLimited      = "RATE_LIMITED"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeUnavailable      = "SERVICE_UNAVAILABLE"
)

const errBaseURL = "https://api.teslapay.eu/errors"

// NewValidationError returns a 400 validation error.
func NewValidationError(detail, instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/validation-error",
		Title:     "Validation Error",
		Status:    http.StatusBadRequest,
		Detail:    detail,
		Instance:  instance,
		ErrorCode: ErrCodeValidation,
	}
}

// NewAuthError returns a 401 authentication error.
func NewAuthError(code, detail string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/authentication-error",
		Title:     "Authentication Error",
		Status:    http.StatusUnauthorized,
		Detail:    detail,
		ErrorCode: code,
	}
}

// NewForbiddenError returns a 403 forbidden error.
func NewForbiddenError(code, detail string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/forbidden",
		Title:     "Forbidden",
		Status:    http.StatusForbidden,
		Detail:    detail,
		ErrorCode: code,
	}
}

// NewNotFoundError returns a 404 not found error.
func NewNotFoundError(detail, instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/not-found",
		Title:     "Not Found",
		Status:    http.StatusNotFound,
		Detail:    detail,
		Instance:  instance,
		ErrorCode: ErrCodeNotFound,
	}
}

// NewConflictError returns a 409 conflict error.
func NewConflictError(detail string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/conflict",
		Title:     "Conflict",
		Status:    http.StatusConflict,
		Detail:    detail,
		ErrorCode: ErrCodeConflict,
	}
}

// NewBusinessError returns a 422 unprocessable entity error for business rule violations.
func NewBusinessError(code, title, detail string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/" + code,
		Title:     title,
		Status:    http.StatusUnprocessableEntity,
		Detail:    detail,
		ErrorCode: code,
	}
}

// NewRateLimitError returns a 429 rate limit error.
func NewRateLimitError() *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/rate-limited",
		Title:     "Rate Limit Exceeded",
		Status:    http.StatusTooManyRequests,
		Detail:    "You have exceeded the request rate limit. Please try again later.",
		ErrorCode: ErrCodeRateLimited,
	}
}

// NewInternalError returns a 500 internal server error.
// The detail parameter should be a generic message; never expose internal state.
func NewInternalError(traceID string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/internal-error",
		Title:     "Internal Server Error",
		Status:    http.StatusInternalServerError,
		Detail:    "An unexpected error occurred. Please contact support if the issue persists.",
		ErrorCode: ErrCodeInternal,
		TraceID:   traceID,
	}
}

// NewUnavailableError returns a 503 service unavailable error.
func NewUnavailableError(detail string) *ProblemDetail {
	return &ProblemDetail{
		Type:      errBaseURL + "/service-unavailable",
		Title:     "Service Unavailable",
		Status:    http.StatusServiceUnavailable,
		Detail:    detail,
		ErrorCode: ErrCodeUnavailable,
	}
}
