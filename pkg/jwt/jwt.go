package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	userID uint   `json:"user_id"`
	email  string `json:"email"`
	role   string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateToken(UserID uint, Email string, Role string) (string, error) {
	Claims := Claims{
		userID: UserID,
		email:  Email,
		role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour),
			IssuedAt:  time.Now().Add(24 * time.Hour),
		},
	}
}
func ValidateToken(TokenString string) (*Claims, error) {

}
