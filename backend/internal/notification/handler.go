package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/teslapay/backend/internal/common"
)

// Handler provides HTTP request handlers for the notification API endpoints.
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new notification handler.
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes sets up the notification routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	notifs := rg.Group("/notifications")
	{
		notifs.GET("", h.ListNotifications)
		notifs.POST("/:id/read", h.MarkAsRead)
		notifs.GET("/preferences", h.GetPreferences)
		notifs.PUT("/preferences", h.UpdatePreferences)
	}
}

// ListNotifications handles GET /notifications
func (h *Handler) ListNotifications(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		problem := common.NewValidationError("user_id not found in context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		problem := common.NewValidationError("invalid user_id in context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid user_id format", c.Request.URL.Path)
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

	notifs, err := h.service.ListNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if notifs == nil {
		notifs = []*Notification{}
	}

	c.JSON(http.StatusOK, ListNotificationsResponse{
		Data:  notifs,
		Total: len(notifs),
	})
}

// MarkAsRead handles POST /notifications/:id/read
func (h *Handler) MarkAsRead(c *gin.Context) {
	notifIDStr := c.Param("id")
	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid notification id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), notifID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPreferences handles GET /notifications/preferences
func (h *Handler) GetPreferences(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		problem := common.NewValidationError("user_id not found in context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		problem := common.NewValidationError("invalid user_id in context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid user_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	prefs, err := h.service.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences handles PUT /notifications/preferences
func (h *Handler) UpdatePreferences(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		problem := common.NewValidationError("user_id not found in context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		problem := common.NewValidationError("invalid user_id in context", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		problem := common.NewValidationError("Invalid user_id format", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	var req PreferencesUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		problem := common.NewValidationError(err.Error(), c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	updatedPrefs, err := h.service.UpdatePreferences(c.Request.Context(), userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedPrefs)
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
