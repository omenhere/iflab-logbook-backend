package utils

import (
    "errors"
    "net/http"
    "os"
    "time"
	"log"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims represents the JWT claims
type Claims struct {
    ID string `json:"id"`
    jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a given user ID
func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "id":  userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ParseToken validates and parses a JWT token, returning the user ID
func ParseToken(tokenString string) (string, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })

    if err != nil {
        return "", err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return "", errors.New("invalid token")
    }

    log.Printf("Parsed user ID from token: %s", claims.ID)

    if claims.ID == "" {
        return "", errors.New("missing user ID in token")
    }

    return claims.ID, nil
}


// SetTokenCookie sets the JWT token in an HTTP-only cookie
func SetTokenCookie(c *gin.Context, token string) {
    cookie := &http.Cookie{
        Name:     "jwt",
        Value:    token,
        Path:     "/",
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
        Secure:   false, // Change to true in production
    }
    http.SetCookie(c.Writer, cookie)
}
