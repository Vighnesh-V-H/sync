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

	if req.ApiKey == "" {
		h.logger.Warn().
			Str("ip", c.ClientIP()).
			Msg("Missing api_key in event request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}
	if req.ID == "" {
		h.logger.Warn().
			Str("api_key", req.ApiKey).
			Str("ip", c.ClientIP()).
			Msg("Missing event ID in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	if req.Payload == nil {
		h.logger.Warn().
			Str("api_key", req.ApiKey).
			Str("event_id", req.ID).
			Str("ip", c.ClientIP()).
			Msg("Missing payload in event request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload is required"})
		return
	}

	h.logger.Info().
		Str("event_id", req.ID).
		Str("api_key", req.ApiKey).
		Str("ip", c.ClientIP()).
		Msg("Processing event request")

	ctx := c.Request.Context()
	res, err := h.svc.AddEvent(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).
			Str("event_id", req.ID).
			Str("api_key", req.ApiKey).
			Str("ip", c.ClientIP()).
			Msg("Failed to add event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("event_id", req.ID).
		Str("api_key", req.ApiKey).
		Str("ip", c.ClientIP()).
		Msg("Event added successfully")
	c.JSON(http.StatusOK, res)
}
