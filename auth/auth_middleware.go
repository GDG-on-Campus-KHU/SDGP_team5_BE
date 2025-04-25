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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		// extract claims -> store in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, exists := claims["user_id"]; exists {
				// c.Set("user", userID)
				c.Set("user", fmt.Sprintf("%v", userID))
			} else if userName, exists := claims["name"]; exists {
				c.Set("user", userName)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User claim not found"})
				c.Abort()
				return
			}
		}

		// proceed to next handler
		c.Next()
	}
}
