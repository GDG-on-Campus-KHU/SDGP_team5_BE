// auth/jwt.go

package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

func GenerateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
