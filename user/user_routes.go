// user/user_handler.go

package user

import (
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler *UserHandler) {
	// users routes
	r.GET("/api/users/me", auth.JWTAuthMiddleware(), userHandler.GetUserInfo)
	r.PATCH("/api/users/me/country", auth.JWTAuthMiddleware(), userHandler.UpdateCountry)
	r.GET("/api/users/info/:member_id", auth.JWTAuthMiddleware(), userHandler.GetUserInfoByID)
	// r.DELETE("/api/users/me", auth.JWTAuthMiddleware(), userHandler.DeleteUser)
}
