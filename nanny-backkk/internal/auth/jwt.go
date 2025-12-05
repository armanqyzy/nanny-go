package auth

import (
	"time"

	"nanny-backend/internal/common/models"
	"nanny-backend/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims — то, что мы кладём внутрь токена
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT генерирует JWT-токен для пользователя
func GenerateJWT(user *models.User) (string, error) {
	cfg := config.Load()
	secret := []byte(cfg.JWTSecret)

	claims := JWTClaims{
		UserID: user.UserID,
		Role:   user.Role,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // токен на 3 дня
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
