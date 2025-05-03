// medical_info/medical_info_routes.go

package medical_info

import (
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
	"github.com/gin-gonic/gin"
)

func RegisterMedicalInfoRoutes(r *gin.Engine, medicalInfoHandler *MedicalInfoHandler) {
	// medical_info routes
	r.POST("/api/medical-info", auth.JWTAuthMiddleware(), medicalInfoHandler.CreateMedicalInfo)
	r.GET("/api/medical-info/:id", auth.JWTAuthMiddleware(), medicalInfoHandler.GetMedicalInfo)
	r.PUT("/api/medical-info", auth.JWTAuthMiddleware(), medicalInfoHandler.UpdateMedicalInfo)
}
