package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Handler provides HTTP handlers for the account API endpoints.
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new account handler.
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes sets up account routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/me", h.GetProfile)
		users.PATCH("/me", h.UpdateProfile)
	}

	accounts := rg.Group("/accounts")
	{
		accounts.GET("", h.ListAccounts)
		accounts.POST("/:account_id/sub-accounts", h.CreateSubAccount)
		accounts.POST("/:account_id/sub-accounts/:sub_account_id/close", h.CloseSubAccount)
		accounts.GET("/:account_id/transactions", h.GetTransactions)
	}

	beneficiaries := rg.Group("/beneficiaries")
	{
		beneficiaries.GET("", h.ListBeneficiaries)
		beneficiaries.POST("", h.CreateBeneficiary)
	}
}

// GetProfile handles GET /api/v1/users/me.
func (h *Handler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, common.NewAuthError(common.ErrCodeAuthMissing, "Not authenticated"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError("Invalid user ID", c.Request.URL.Path))
		return
	}

	user, err := h.service.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, common.NewNotFoundError("User profile not found", c.Request.URL.Path))
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles PATCH /api/v1/users/me.
func (h *Handler) UpdateProfile(c *gin.Context) {
	// In production: parse partial update fields, validate, apply, check if re-verification needed.
	c.JSON(http.StatusOK, gin.H{"message": "profile update not yet implemented"})
}

// ListAccounts handles GET /api/v1/accounts.
func (h *Handler) ListAccounts(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, common.NewAuthError(common.ErrCodeAuthMissing, "Not authenticated"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError("Invalid user ID", c.Request.URL.Path))
		return
	}

	accounts, err := h.service.GetAccountsByUserID(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var responses []AccountResponse
	for _, acct := range accounts {
		subAccounts, err := h.service.repo.GetSubAccountsByAccountID(c.Request.Context(), acct.ID)
		if err != nil {
			h.handleError(c, err)
			return
		}

		var subResponses []SubAccountResponse
		for _, sa := range subAccounts {
			subResponses = append(subResponses, SubAccountResponse{
				ID:        sa.ID.String(),
				Currency:  sa.Currency,
				IBAN:      sa.IBAN,
				BIC:       sa.BIC,
				Balance: SubAccountBalance{
					Available: "0.00", // In production: fetch from ledger service.
					Pending:   "0.00",
					Total:     "0.00",
				},
				IsDefault: sa.IsDefault,
			})
		}

		responses = append(responses, AccountResponse{
			ID:              acct.ID.String(),
			AccountNumber:   acct.AccountNumber,
			Status:          acct.Status,
			SubAccounts:     subResponses,
			TotalBalanceEUR: "0.00",
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// CreateSubAccount handles POST /api/v1/accounts/:account_id/sub-accounts.
func (h *Handler) CreateSubAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError("Invalid account ID", c.Request.URL.Path))
		return
	}

	var req CreateSubAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError(err.Error(), c.Request.URL.Path))
		return
	}

	subAccount, err := h.service.OpenSubAccount(c.Request.Context(), accountID, req.Currency)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, SubAccountResponse{
		ID:       subAccount.ID.String(),
		Currency: subAccount.Currency,
		IBAN:     subAccount.IBAN,
		BIC:      subAccount.BIC,
		Balance: SubAccountBalance{
			Available: "0.00",
			Pending:   "0.00",
			Total:     "0.00",
		},
		IsDefault: subAccount.IsDefault,
	})
}

// CloseSubAccount handles POST /api/v1/accounts/:account_id/sub-accounts/:sub_account_id/close.
func (h *Handler) CloseSubAccount(c *gin.Context) {
	// In production: validate zero balance, transfer remaining, close.
	c.JSON(http.StatusOK, gin.H{"status": "closed", "remaining_transferred": "0.00"})
}

// GetTransactions handles GET /api/v1/accounts/:account_id/transactions.
func (h *Handler) GetTransactions(c *gin.Context) {
	// In production: query ledger service for transaction history.
	c.JSON(http.StatusOK, gin.H{
		"data":       []any{},
		"pagination": common.PaginationResponse{HasMore: false},
	})
}

// ListBeneficiaries handles GET /api/v1/beneficiaries.
func (h *Handler) ListBeneficiaries(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []any{}})
}

// CreateBeneficiary handles POST /api/v1/beneficiaries.
func (h *Handler) CreateBeneficiary(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "beneficiary created"})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	if problem, ok := err.(*common.ProblemDetail); ok {
		problem.Instance = c.Request.URL.Path
		problem.TraceID = common.RequestIDFromContext(c.Request.Context())
		c.JSON(problem.Status, problem)
		return
	}

	requestID := common.RequestIDFromContext(c.Request.Context())
	h.logger.Error("unexpected error",
		zap.Error(err),
		zap.String("request_id", requestID),
	)
	c.JSON(http.StatusInternalServerError, common.NewInternalError(requestID))
}
