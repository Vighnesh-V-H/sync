package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var validate = validator.New()

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func SetupAuthRoutes(router *gin.Engine, log zerolog.Logger) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", signupHandler(log))
		auth.POST("/signin", signinHandler(log))
	}
}

func signupHandler(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignupRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error().Err(err).Msg("Failed to bind signup request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Error().Err(err).Msg("Signup validation failed")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		log.Info().
			Str("email", req.Email).
			Str("name", req.Name).
			Msg("Signup request received and validated")

		c.JSON(http.StatusOK, gin.H{
			"message": "Signup request validated successfully",
			"email":   req.Email,
			"name":    req.Name,
		})
	}
}

func signinHandler(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SigninRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error().Err(err).Msg("Failed to bind signin request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Error().Err(err).Msg("Signin validation failed")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		log.Info().
			Str("email", req.Email).
			Msg("Signin request received and validated")

		c.JSON(http.StatusOK, gin.H{
			"message": "Signin request validated successfully",
			"email":   req.Email,
		})
	}
}
