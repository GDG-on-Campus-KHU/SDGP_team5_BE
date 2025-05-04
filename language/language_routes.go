// language/language_routes.go

package language

import (
	"github.com/gin-gonic/gin"
	
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
)

func RegisterLanguageRoutes(r *gin.Engine, translationHandler *TranslationHandler) {

	// medical information translation routes
	r.POST("/api/medical-info/translate", auth.JWTAuthMiddleware(), translationHandler.TranslateMedicalInfoHandler)
}