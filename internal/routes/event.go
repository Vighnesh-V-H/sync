package routes

import (
	"github.com/Vighnesh-V-H/sync/internal/handler"
	"github.com/Vighnesh-V-H/sync/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(router gin.IRouter, h *handler.EventHandler, secret string) {
	event := router.Group("/event")
	event.Use(middleware.AuthMiddleware(secret))
	{
		event.POST("/add", h.AddEvent)
	}
}
