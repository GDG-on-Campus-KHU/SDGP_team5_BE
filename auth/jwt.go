// auth/jwt.go

package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)


func generateToken(userID, name, email string, expiration time.Duration) (string, error) {
    claims := &CustomClaims{
        RegisteredClaims: &jwt.RegisteredClaims{
            Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
        },
        Name:  name,
        Email: email,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secretKey := os.Getenv("JWT_SECRET_KEY")

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }
    return tokenString, nil
}


func GenerateAccessToken(userID, name, email string, expiration time.Duration) (string, error) {
    return generateToken(userID, name, email, expiration)
}


func GenerateRefreshToken(userID, name, email string) (string, error) {
    return generateToken(userID, name, email, 24*time.Hour)
}