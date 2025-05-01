// recording/recording_handler.go

package recording

import (
	"net/http"
	"strconv"
	"fmt"
	"os"
	"time"
	"io"

	"github.com/gin-gonic/gin"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	coreUtil "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/recording/util"
)


// POST /api/recordings/stt
// @Summary      Transcribe short audio
// @Description  Accepts a short audio file and returns the transcribed text using Speech-to-Text (STT).
// @Tags         Recordings
// @Accept       multipart/form-data
// @Produce      json
// @Param        audio  formData  file  true  "Short audio file for transcription"
// @Success      200    {object}  map[string]string  "Returns the transcription result"
// @Failure      400    {object}  map[string]string  "Bad Request"
// @Failure      500    {object}  map[string]string  "Internal Server Error"
// @Router       /api/recordings/stt [post]
// @Security     BearerAuth
func SyncSttRecordingHandler(c *gin.Context) {

    ctx := c.Request.Context()

    // extract user ID from Gin context
    userIDStr, err := coreUtil.GetUserIDFromContext(c)
    if err != nil {
        coreUtil.RespondUnauthorized(c, "Unauthorized access")
        return
    }

    // convert userID to int
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        coreUtil.RespondBadRequest(c, "Invalid user ID format")
        return
    }

    // get user info by userID
    user, err := coreUtil.GetUserByUserID(ctx, userID)
    if err != nil {
        coreUtil.RespondInternalError(c, "Failed to get user info")
        return
    }

    // language code from user info (app_lang)
	languageCode := coreUtil.AppLangToLangCode(user.AppLang)
    fmt.Printf("Language Code: %s\n", languageCode)

	// get uploaded audio file
	fileHeader, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
		return
	}

	// temporary local copy: *multipart.FileHeader to *os.File
	srcFile, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer srcFile.Close()

	// temporary file to pass to STT function
	tempFilePath := fmt.Sprintf("tmp/%d.wav", time.Now().Unix())
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file"})
		return
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, srcFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy uploaded file"})
		return
	}

	// call STT service
	transcription, err := util.SynchronousSpeechToText(tempFilePath, languageCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transcribe audio: %v", err)})
		return
	}

	// upload file to GCS
	recordingURL, err := util.UploadFileToGCS(fileHeader, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to GCS"})
		return
	}

	// create recording record
	record := NewRecording(userID, recordingURL)
	record.AppLang = string(user.AppLang)
	record.RecordingText = transcription

	// insert into database
	_, err = dbConfig.RecordingCollection.InsertOne(ctx, record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recording to DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcription":  transcription,
		"recording_url":  recordingURL,
	})
}



// POST /api/recordings/upload
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
// @Router /api/recordings/upload [post]
// @Security BearerAuth

// func CreateRecording(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	// extract user ID from Gin context
// 	userIDStr, err := coreUtil.GetUserIDFromContext(c)
// 	if err != nil {
// 		coreUtil.RespondUnauthorized(c, "Unauthorized access")
// 		return
// 	}

// 	userID, err := strconv.Atoi(userIDStr)
// 	if err != nil {
// 		coreUtil.RespondBadRequest(c, "Invalid user ID format")
// 		return
// 	}

// 	user, err := coreUtil.GetUserByUserID(ctx, userID)
// 	if err != nil {
// 		coreUtil.RespondInternalError(c, "Failed to get user info")
// 		return
// 	}

// 	countryCode := user.CountryCode
// 	if countryCode == "" {
// 		coreUtil.RespondBadRequest(c, "Country code not found for user")
// 		return
// 	}

// 	// multipart file
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "No file is attached"})
// 		return
// 	}

// 	// save audio recording file with metadata
// 	recording, err := SaveRecording(ctx, userID, countryCode, file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"recording_id":  recording.RecordingID,
// 		"recording_url": recording.RecordingURL,
// 		"message":       "Recording uploaded. Transcription will be processed.",
// 	})
// }



// POST /api/recordings
// CreateRecordingHandler godoc
// @Summary      Create a new recording
// @Description  Uploads an audio file and creates a new recording. The server stores the file and returns the recording URL.
// @Tags         Recordings
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Audio file to upload"
// @Success      200   {object}  map[string]string  "recording_url Returned"
// @Failure      400   {object}  map[string]string  "Bad Request"
// @Failure      401   {object}  map[string]string  "Unauthorized"
// @Failure      500   {object}  map[string]string  "Internal Server Error"
// @Router       /api/recordings [post]
// @Security BearerAuth
func CreateRecordingHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := coreUtil.GetUserIDFromContext(c)
	if err != nil {
		coreUtil.RespondUnauthorized(c, "Unauthorized: user ID not found")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		coreUtil.RespondBadRequest(c, "Invalid user ID")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		coreUtil.RespondBadRequest(c, "File is required")
		return
	}

	recordingURL, err := CreateRecordingService(ctx, userID, file)
	if err != nil {
		coreUtil.RespondInternalError(c, err.Error())
		return
	}

	coreUtil.RespondSuccess(c, gin.H{
		"recording_url": recordingURL,
	})
}