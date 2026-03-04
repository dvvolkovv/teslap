package card

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Handler provides HTTP request handlers for the card API endpoints.
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new card handler.
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes sets up the card routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	cards := rg.Group("/cards")
	{
		cards.POST("/virtual", h.IssueVirtualCard)
		cards.POST("/physical", h.IssuePhysicalCard)
		cards.GET("", h.ListCards)
		cards.GET("/:card_id", h.GetCard)
		cards.POST("/:card_id/activate", h.ActivateCard)
		cards.POST("/:card_id/freeze", h.FreezeCard)
		cards.POST("/:card_id/unfreeze", h.UnfreezeCard)
		cards.POST("/:card_id/block", h.BlockCard)
		cards.PUT("/:card_id/controls", h.UpdateCardControls)
		cards.GET("/:card_id/transactions", h.GetCardTransactions)
	}
}

// IssueVirtualCard handles POST /cards/virtual.
func (h *Handler) IssueVirtualCard(c *gin.Context) {
	var req IssueCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.IssueVirtualCard(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toCardResponse(card))
}

// IssuePhysicalCard handles POST /cards/physical.
func (h *Handler) IssuePhysicalCard(c *gin.Context) {
	var req IssueCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.IssuePhysicalCard(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toCardResponse(card))
}

// GetCard handles GET /cards/:card_id.
func (h *Handler) GetCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.GetCard(c.Request.Context(), cardID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCardResponse(card))
}

// ListCards handles GET /cards?account_id=...
func (h *Handler) ListCards(c *gin.Context) {
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

	cards, err := h.service.ListCards(c.Request.Context(), accountID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var responses []*CardResponse
	for _, card := range cards {
		responses = append(responses, toCardResponse(card))
	}
	if responses == nil {
		responses = []*CardResponse{}
	}

	c.JSON(http.StatusOK, ListCardsResponse{
		Data:  responses,
		Total: len(responses),
	})
}

// ActivateCard handles POST /cards/:card_id/activate.
func (h *Handler) ActivateCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	var req ActivateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.ActivateCard(c.Request.Context(), cardID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCardResponse(card))
}

// FreezeCard handles POST /cards/:card_id/freeze.
func (h *Handler) FreezeCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.FreezeCard(c.Request.Context(), cardID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCardResponse(card))
}

// UnfreezeCard handles POST /cards/:card_id/unfreeze.
func (h *Handler) UnfreezeCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.UnfreezeCard(c.Request.Context(), cardID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCardResponse(card))
}

// BlockCard handles POST /cards/:card_id/block.
func (h *Handler) BlockCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.BlockCard(c.Request.Context(), cardID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCardResponse(card))
}

// UpdateCardControls handles PUT /cards/:card_id/controls.
func (h *Handler) UpdateCardControls(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	var req CardControlsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	card, err := h.service.UpdateControls(c.Request.Context(), cardID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCardResponse(card))
}

// GetCardTransactions handles GET /cards/:card_id/transactions?limit=...&offset=...
func (h *Handler) GetCardTransactions(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("card_id"))
	if err != nil {
		problem := common.NewValidationError("Invalid card_id format", c.Request.URL.Path)
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

	txs, err := h.service.GetCardTransactions(c.Request.Context(), cardID, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var responses []*CardTransactionResponse
	for _, tx := range txs {
		responses = append(responses, toCardTransactionResponse(tx))
	}
	if responses == nil {
		responses = []*CardTransactionResponse{}
	}

	c.JSON(http.StatusOK, ListCardTransactionsResponse{
		Data:  responses,
		Total: len(responses),
	})
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
