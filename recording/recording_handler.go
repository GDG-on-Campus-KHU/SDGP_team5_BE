// recording/recording_handler.go

package recording

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)

// POST /api/recordings
// @Summary Upload a new recording
// @Description Upload an audio recording for a user. The file will be saved and processed for transcription.
// @Tags recordings
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Audio file to be uploaded"
// @Success 200 {object} map[string]interface{} "Recording uploaded successfully"
// @Failure 400 {object} map[string]string "Bad request, invalid file or user"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/recordings [post]
func CreateRecording(c *gin.Context) {
	ctx := c.Request.Context()

	// extract user ID from Gin context
	userIDStr, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID format")
		return
	}

	user, err := util.GetUserByUserID(ctx, userID)
	if err != nil {
		util.RespondInternalError(c, "Failed to get user info")
		return
	}

	countryCode := user.CountryCode
	if countryCode == "" {
		util.RespondBadRequest(c, "Country code not found for user")
		return
	}

	// multipart file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No file is attached"})
		return
	}

	// save audio recording file with metadata
	recording, err := SaveRecording(ctx, userID, countryCode, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recording_id":  recording.RecordingID,
		"recording_url": recording.RecordingURL,
		"message":       "Recording uploaded. Transcription will be processed.",
	})
}