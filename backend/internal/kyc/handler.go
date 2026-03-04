package kyc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Handler provides HTTP request handlers for the KYC API endpoints.
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new KYC handler.
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes sets up the authenticated KYC routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	kyc := rg.Group("/kyc")
	{
		kyc.POST("/start", h.StartVerification)
		kyc.GET("/status", h.GetStatus)
	}
}

// RegisterWebhookRoute registers the public (unauthenticated) webhook route.
func (h *Handler) RegisterWebhookRoute(rg *gin.RouterGroup) {
	rg.POST("/kyc/webhook", h.HandleWebhook)
}

// StartVerification handles POST /kyc/start.
func (h *Handler) StartVerification(c *gin.Context) {
	var req KYCStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	// Override user_id from JWT context if present.
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, ok := userIDStr.(string); ok && uid != "" {
			req.UserID = uid
		}
	}

	resp, err := h.service.StartVerification(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetStatus handles GET /kyc/status.
func (h *Handler) GetStatus(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		problem := common.NewValidationError("user_id is required in JWT context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	uid, ok := userIDStr.(string)
	if !ok || uid == "" {
		problem := common.NewValidationError("user_id is required in JWT context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	userID, err := uuid.Parse(uid)
	if err != nil {
		problem := common.NewValidationError("Invalid user_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	resp, err := h.service.GetStatus(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleWebhook handles POST /kyc/webhook (public, unauthenticated).
func (h *Handler) HandleWebhook(c *gin.Context) {
	var payload KYCWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	if err := h.service.HandleWebhook(c.Request.Context(), &payload); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
