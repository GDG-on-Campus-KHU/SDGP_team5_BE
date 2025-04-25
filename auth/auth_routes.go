// auth/auth_routes.go

package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.GET("/api/auth/login", LoginHandler)
	r.GET("/api/auth/callback", CallbackHandler)
	r.GET("/api/auth/protected", JWTAuthMiddleware(), ProtectedHandler)
	r.POST("/api/auth/refresh-token", RefreshTokenHandler)
}