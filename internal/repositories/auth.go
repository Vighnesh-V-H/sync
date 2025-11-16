package repositories

import (
	"context"
	"time"

	"github.com/Vighnesh-V-H/sync/internal/db"
	errors "github.com/Vighnesh-V-H/sync/internal/error"
	"github.com/Vighnesh-V-H/sync/internal/models"
	"github.com/Vighnesh-V-H/sync/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
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
	mail := user.Email

	existing, err := r.GetUserByEmail(ctx, mail)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.ErrUserAlreadyExists
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	apikey, err := gonanoid.New()
	if err != nil {
		return err
	}
	syncID := "sync_" + apikey

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, email, password, name, api_key, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, user.ID, user.Email, user.Password, user.Name, syncID, true, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, api_key, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Api_Key,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
