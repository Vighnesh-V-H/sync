package routes

import (
	"github.com/Vighnesh-V-H/sync/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router gin.IRouter, h *handler.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.Signup)
		auth.POST("/signin", h.Signin)
	}
}
