package routes

import (
	"github.com/Vighnesh-V-H/sync/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(router gin.IRouter, h *handler.EventHandler) {
	event := router.Group("/event")
	{
		event.POST("", h.AddEvent)
	}
}
