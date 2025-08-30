package api

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(userID int64, ttl time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
    
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(ttl).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}