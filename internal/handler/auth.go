package handler

import (
	"net/http"

	"github.com/Vighnesh-V-H/sync/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var validate = validator.New()

type AuthHandler struct {
	svc    *service.AuthService
	logger zerolog.Logger
}

func NewAuthHandler(svc *service.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		svc:    svc,
		logger: logger.With().Str("handler", "auth").Logger(),
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req service.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Msg("Failed to bind signup request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		h.logger.Warn().Err(err).
			Str("email", req.Email).
			Str("ip", c.ClientIP()).
			Msg("Signup validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("email", req.Email).
		Str("name", req.Name).
		Str("ip", c.ClientIP()).
		Msg("Attempting user signup")

	ctx := c.Request.Context()
	res, err := h.svc.Signup(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).
			Str("email", req.Email).
			Str("ip", c.ClientIP()).
			Msg("Signup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("email", req.Email).
		Str("ip", c.ClientIP()).
		Msg("User signup successful")
	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Signin(c *gin.Context) {
	var req service.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Msg("Failed to bind signin request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		h.logger.Warn().Err(err).
			Str("email", req.Email).
			Str("ip", c.ClientIP()).
			Msg("Signin validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("email", req.Email).
		Str("ip", c.ClientIP()).
		Msg("Attempting user signin")

	ctx := c.Request.Context()
	res, err := h.svc.Signin(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).
			Str("email", req.Email).
			Str("ip", c.ClientIP()).
			Msg("Signin failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("email", req.Email).
		Str("ip", c.ClientIP()).
		Msg("User signin successful")
	c.JSON(http.StatusOK, res)
}
