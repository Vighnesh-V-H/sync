package repositories

import (
	"context"

	"github.com/Vighnesh-V-H/sync/internal/db"
	"github.com/Vighnesh-V-H/sync/internal/models"
	"github.com/rs/zerolog"
)

type AuthRepository struct {
	db  *db.DB
	log zerolog.Logger
}

func NewAuthRepository(db *db.DB, log zerolog.Logger) *AuthRepository {
	return &AuthRepository{
		db:  db,
		log: log,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.log.Info().Msg("creating user")
	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.log.Info().Str("email", email).Msg("getting user by email")
	return &models.User{
		Email: email,
	}, nil
}
