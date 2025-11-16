package service

import (
	"context"
	"errors"

	internalErrors "github.com/Vighnesh-V-H/sync/internal/error"
	"github.com/Vighnesh-V-H/sync/internal/models"
	"github.com/Vighnesh-V-H/sync/internal/repositories"
	"github.com/Vighnesh-V-H/sync/internal/utils"
)

type AuthService struct {
	repo      *repositories.AuthRepository
	jwtConfig utils.JWTConfig
}

func NewAuthService(repo *repositories.AuthRepository, jwtCfg utils.JWTConfig) *AuthService {
	return &AuthService{repo: repo, jwtConfig: jwtCfg}
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
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &SignupResponse{
		Success: true,
		Message: "Signup successful. Please verify your email.",
	}, nil
}

func (s *AuthService) Signin(ctx context.Context, req SigninRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, internalErrors.ErrUserNotFound
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}
