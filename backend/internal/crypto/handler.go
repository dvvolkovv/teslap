package crypto

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Handler provides HTTP request handlers for the crypto API endpoints.
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new crypto handler.
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes sets up the crypto routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	crypto := rg.Group("/crypto")
	{
		crypto.GET("/wallet", h.GetWallet)
		crypto.GET("/prices", h.GetPrices)
		crypto.GET("/quote", h.GetQuote)
		crypto.POST("/buy", h.BuyCrypto)
		crypto.POST("/sell", h.SellCrypto)
		crypto.POST("/send", h.SendCrypto)
		crypto.GET("/transactions", h.GetTransactions)
	}
}

// ExecuteTradeRequest is the request body for POST /crypto/buy and /crypto/sell.
type ExecuteTradeRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	QuoteID string `json:"quote_id" binding:"required"`
}

// SendCryptoHandlerRequest is the request body for POST /crypto/send.
type SendCryptoHandlerRequest struct {
	UserID           string `json:"user_id" binding:"required"`
	Symbol           string `json:"symbol" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
	RecipientAddress string `json:"recipient_address" binding:"required"`
}

// GetWallet handles GET /crypto/wallet?user_id=...
func (h *Handler) GetWallet(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		problem := common.NewValidationError("user_id query parameter is required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	wallet, balances, err := h.service.GetOrCreateWallet(c.Request.Context(), userIDStr)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toWalletResponse(wallet, balances, cryptoPrices))
}

// GetPrices handles GET /crypto/prices.
func (h *Handler) GetPrices(c *gin.Context) {
	prices, err := h.service.GetPrices(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": prices})
}

// GetQuote handles GET /crypto/quote?user_id=...&action=buy&symbol=FUSE&amount=100
func (h *Handler) GetQuote(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		problem := common.NewValidationError("user_id query parameter is required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	var req GetQuoteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	quote, err := h.service.GetQuote(c.Request.Context(), userIDStr, req.Action, req.Symbol, req.Amount)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, quote)
}

// BuyCrypto handles POST /crypto/buy.
func (h *Handler) BuyCrypto(c *gin.Context) {
	var req ExecuteTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	tx, err := h.service.BuyCrypto(c.Request.Context(), req.UserID, req.QuoteID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTransactionResponse(tx))
}

// SellCrypto handles POST /crypto/sell.
func (h *Handler) SellCrypto(c *gin.Context) {
	var req ExecuteTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	tx, err := h.service.SellCrypto(c.Request.Context(), req.UserID, req.QuoteID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTransactionResponse(tx))
}

// SendCrypto handles POST /crypto/send.
func (h *Handler) SendCrypto(c *gin.Context) {
	var req SendCryptoHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	tx, err := h.service.SendCrypto(c.Request.Context(), req.UserID, req.Symbol, req.Amount, req.RecipientAddress)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTransactionResponse(tx))
}

// GetTransactions handles GET /crypto/transactions?user_id=...&limit=20&offset=0
func (h *Handler) GetTransactions(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		problem := common.NewValidationError("user_id query parameter is required", c.Request.URL.Path)
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

	txs, err := h.service.GetTransactions(c.Request.Context(), userIDStr, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]*TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		responses = append(responses, toTransactionResponse(tx))
	}

	c.JSON(http.StatusOK, ListTransactionsResponse{
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
