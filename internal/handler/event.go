package handler

import (
	"net/http"

	"github.com/Vighnesh-V-H/sync/internal/service"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) AddEvent(c *gin.Context) {
	var req service.AddEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}


	if req.ApiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}
	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	if req.Payload == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload is required"})
		return
	}

	ctx := c.Request.Context()
	res, err := h.svc.AddEvent(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
