package service

import (
	"context"
	"errors"

	internalErrors "github.com/Vighnesh-V-H/sync/internal/error"
	"github.com/Vighnesh-V-H/sync/internal/models"
	"github.com/Vighnesh-V-H/sync/internal/repositories"
	"github.com/Vighnesh-V-H/sync/internal/utils"
	"github.com/rs/zerolog"
)

type AuthService struct {
	repo      *repositories.AuthRepository
	jwtConfig utils.JWTConfig
	logger    zerolog.Logger
}

func NewAuthService(repo *repositories.AuthRepository, jwtCfg utils.JWTConfig, logger zerolog.Logger) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtConfig: jwtCfg,
		logger:    logger.With().Str("service", "auth").Logger(),
	}
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
	s.logger.Debug().
		Str("email", req.Email).
		Str("name", req.Name).
		Msg("Starting user signup process")

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.Error().Err(err).
			Str("email", req.Email).
			Msg("Failed to create user in repository")
		return nil, err
	}

	s.logger.Info().
		Str("email", req.Email).
		Str("user_id", user.ID).
		Msg("User created successfully")

	return &SignupResponse{
		Success: true,
		Message: "Signup successful. Please verify your email.",
	}, nil
}

func (s *AuthService) Signin(ctx context.Context, req SigninRequest) (*AuthResponse, error) {
	s.logger.Debug().
		Str("email", req.Email).
		Msg("Starting user signin process")

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error().Err(err).
			Str("email", req.Email).
			Msg("Failed to fetch user from repository")
		return nil, err
	}

	if user == nil {
		s.logger.Warn().
			Str("email", req.Email).
			Msg("User not found during signin attempt")
		return nil, internalErrors.ErrUserNotFound
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		s.logger.Warn().
			Str("email", req.Email).
			Str("user_id", user.ID).
			Msg("Invalid password attempt")
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, s.jwtConfig)
	if err != nil {
		s.logger.Error().Err(err).
			Str("email", req.Email).
			Str("user_id", user.ID).
			Msg("Failed to generate JWT token")
		return nil, err
	}

	s.logger.Info().
		Str("email", req.Email).
		Str("user_id", user.ID).
		Msg("User signin successful")

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}
