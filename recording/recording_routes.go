// recording/recording_routes.go

package recording

import (
	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
)

func RegisterRecordingRoutes(r *gin.Engine) {
	// r.POST("/api/recordings/upload", auth.JWTAuthMiddleware(), CreateRecording)
	// r.POST("/api/recordings", auth.JWTAuthMiddleware(), CreateRecordingHandler)
	r.POST("/api/recordings/stt", auth.JWTAuthMiddleware(), SyncSttRecordingHandler)
}