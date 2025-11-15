package service

import (
	"context"

	"github.com/Vighnesh-V-H/sync/internal/models"
	"github.com/Vighnesh-V-H/sync/internal/repositories"
)

type AuthService struct {
	repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (*models.User, error) {
	user := &models.User{
		Email:    req.Email,
		Password: req.Password, 
		Name:     req.Name,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Signin(ctx context.Context, req SigninRequest) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, req.Email)
}
