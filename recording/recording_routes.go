// recording/recording_routes.go

package recording

import (
	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
)

func RegisterRecordingRoutes(r *gin.Engine) {
	recordingService := &RecordingService{}
	recordingHandler := NewRecordingHandler(recordingService)

	r.POST("/api/recordings/stt", auth.JWTAuthMiddleware(), SyncSttRecordingHandler)
	r.GET("/api/recordings/me", auth.JWTAuthMiddleware(), recordingHandler.GetMyRecordings)
	r.GET("/api/recordings/:id", auth.JWTAuthMiddleware(), recordingHandler.GetRecordingByID)
}