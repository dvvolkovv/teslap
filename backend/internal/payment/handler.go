package payment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Handler provides HTTP request handlers for the payment API endpoints.
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new payment handler.
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes sets up the payment routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	payments := rg.Group("/payments")
	{
		payments.POST("/internal", h.CreateInternalPayment)
		payments.POST("/sepa", h.CreateSEPAPayment)
		payments.GET("", h.ListPayments)
		payments.GET("/fx/quote", h.GetFXQuote)
		payments.POST("/fx/execute", h.ExecuteFX)
		payments.POST("/scheduled", h.CreateScheduledPayment)
		payments.GET("/scheduled", h.ListScheduledPayments)
		payments.GET("/:payment_id", h.GetPayment)
	}
}

// CreateInternalPayment handles POST /payments/internal.
func (h *Handler) CreateInternalPayment(c *gin.Context) {
	var req InternalPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	payment, err := h.service.CreateInternalPayment(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPaymentResponse(payment))
}

// CreateSEPAPayment handles POST /payments/sepa.
func (h *Handler) CreateSEPAPayment(c *gin.Context) {
	var req SEPAPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	payment, err := h.service.CreateSEPAPayment(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPaymentResponse(payment))
}

// GetPayment handles GET /payments/:payment_id.
func (h *Handler) GetPayment(c *gin.Context) {
	paymentIDStr := c.Param("payment_id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid payment_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	payment, err := h.service.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPaymentResponse(payment))
}

// ListPayments handles GET /payments?account_id=...&limit=...&offset=...
func (h *Handler) ListPayments(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		problem := common.NewValidationError("account_id query parameter is required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid account_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	payments, err := h.service.ListPayments(c.Request.Context(), accountID, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var responses []*PaymentResponse
	for _, p := range payments {
		responses = append(responses, toPaymentResponse(p))
	}
	if responses == nil {
		responses = []*PaymentResponse{}
	}

	c.JSON(http.StatusOK, ListPaymentsResponse{
		Data:  responses,
		Total: len(responses),
	})
}

// GetFXQuote handles GET /payments/fx/quote?from=EUR&to=USD&amount=100
func (h *Handler) GetFXQuote(c *gin.Context) {
	var req FXQuoteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	quote, err := h.service.GetFXQuote(c.Request.Context(), req.From, req.To, req.Amount)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, quote)
}

// ExecuteFX handles POST /payments/fx/execute.
func (h *Handler) ExecuteFX(c *gin.Context) {
	var req FXExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	if err := h.service.ExecuteFX(c.Request.Context(), &req); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted", "quote_id": req.QuoteID})
}

// CreateScheduledPayment handles POST /payments/scheduled.
func (h *Handler) CreateScheduledPayment(c *gin.Context) {
	var req ScheduledPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	sp, err := h.service.CreateScheduledPayment(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toScheduledPaymentResponse(sp))
}

// ListScheduledPayments handles GET /payments/scheduled?account_id=...
func (h *Handler) ListScheduledPayments(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		problem := common.NewValidationError("account_id query parameter is required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid account_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	scheduled, err := h.service.ListScheduledPayments(c.Request.Context(), accountID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var responses []*ScheduledPaymentResponse
	for _, sp := range scheduled {
		responses = append(responses, toScheduledPaymentResponse(sp))
	}
	if responses == nil {
		responses = []*ScheduledPaymentResponse{}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// handleError converts domain errors to RFC 7807 JSON responses.
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
		zap.String("path", c.Request.URL.Path),
	)
	problem := common.NewInternalError(requestID)
	c.JSON(http.StatusInternalServerError, problem)
}
