package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

func GenerateJWT(id string, email string, apiKey string, cfg JWTConfig) (string, error) {
	claims := jwt.MapClaims{
		"id":      id,
		"email":   email,
		"api_key": apiKey,
		"exp":     time.Now().Add(cfg.Expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
