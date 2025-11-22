package handler

import (
	"net/http"

	"github.com/Vighnesh-V-H/sync/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type EventHandler struct {
	svc    *service.EventService
	logger zerolog.Logger
}

func NewEventHandler(svc *service.EventService, logger zerolog.Logger) *EventHandler {
	return &EventHandler{
		svc:    svc,
		logger: logger.With().Str("handler", "event").Logger(),
	}
}

func (h *EventHandler) AddEvent(c *gin.Context) {
	apiKeyValue, exists := c.Get("api_key")
	if !exists {
		h.logger.Error().
			Str("ip", c.ClientIP()).
			Msg("API key not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	apiKey, ok := apiKeyValue.(string)
	if !ok {
		h.logger.Error().
			Str("ip", c.ClientIP()).
			Msg("Invalid API key type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var req service.AddEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Msg("Failed to bind event request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Payload == nil {
		h.logger.Warn().
			Str("api_key", apiKey).
			Str("ip", c.ClientIP()).
			Msg("Missing payload in event request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload is required"})
		return
	}

	h.logger.Info().
		Str("api_key", apiKey).
		Str("ip", c.ClientIP()).
		Msg("Processing event request")

	ctx := c.Request.Context()
	res, err := h.svc.AddEvent(ctx, apiKey, req)
	if err != nil {
		h.logger.Error().Err(err).
			Str("api_key", apiKey).
			Str("ip", c.ClientIP()).
			Msg("Failed to add event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	h.logger.Info().
		Str("event_id", res.EventID).
		Str("api_key", apiKey).
		Str("ip", c.ClientIP()).
		Msg("Event added successfully")
	
	c.JSON(http.StatusAccepted, res)

}
