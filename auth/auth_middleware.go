// auth/auth_middleware.go

package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)


type CustomClaims struct {
	*jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}


type RefreshToken struct {
	Token string `json:"token"`
}


func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		// token paring & validation
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := os.Getenv("JWT_SECRET_KEY")

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// HMAC: Hash-based Message Authentication Code
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		// exception handling (invalid token)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// extract claims
			if claims, ok := token.Claims.(*CustomClaims); ok {
				c.Set("user", claims.Name)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User claim not found"})
				c.Abort()
				return
			}	

		// proceed to next handler
		c.Next()
	}
}