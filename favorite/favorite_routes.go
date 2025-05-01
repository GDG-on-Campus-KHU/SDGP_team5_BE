// favorite/favorite_handler.go

package favorite

import (
	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
)

func RegisterFavoriteRoutes(r *gin.Engine, favoriteHandler *FavoriteHandler) {
	// favorite routes
	r.POST("/api/favorites/:index", auth.JWTAuthMiddleware(), favoriteHandler.AddFavorite)
	// r.GET("/api/favorites", auth.JWTAuthMiddleware(), favoriteHandler.GetFavorites)
	// r.DELETE("/api/favorites/:id", auth.JWTAuthMiddleware(), favoriteHandler.DeleteFavorite)
}
