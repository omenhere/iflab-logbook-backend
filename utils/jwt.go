package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(nim string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nim": nim,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}
