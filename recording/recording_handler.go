// recording/recording_handler.go

package recording

import (
	"net/http"
	"strconv"
	"fmt"
	"os"
	"time"
	"io"
	"log"
	"path/filepath"

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


	// save with original file extension
	err = os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory"})
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
    originalPath := fmt.Sprintf("tmp/%d%s", time.Now().Unix(), ext)

    outFile, err := os.Create(originalPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file"})
        return
    }
    defer outFile.Close()


	// copy to temporary file
	_, err = io.Copy(outFile, srcFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy uploaded file"})
		return
	}

    // convert to 'wav' format
    convertedPath, err := util.AudioFileConvert(originalPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Audio conversion failed"})
        return
    }

	// call STT service
	transcription, err := util.SynchronousSpeechToText(convertedPath, languageCode)
	if err != nil {
		log.Printf("Error transcribing audio: %v", err)
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

	coreUtil.RespondSuccess(c, gin.H{
		"transcription": transcription,
		"recording_url": recordingURL,
	})
}


type RecordingHandler struct {
	service *RecordingService
}

func NewRecordingHandler(service *RecordingService) *RecordingHandler {
	return &RecordingHandler{
		service: service,
	}
}

// GET /api/recordings/me
// @Summary Get all recordings for the logged-in user
// @Description Retrieve the recordings of the logged-in user
// @Tags recordings
// @Accept json
// @Produce json
// @Success 200 {array} model.Recording "List of recordings the user owns"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/recordings/me [get]
func (h *RecordingHandler) GetMyRecordings(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := coreUtil.GetUserIDFromContext(c)
	if err != nil {
		coreUtil.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		coreUtil.RespondUnauthorized(c, "Invalid user ID")
		return
	}

	recordings, err := h.service.GetRecordingsByUserID(ctx, userID)
	if err != nil {
		coreUtil.RespondInternalError(c, "Failed to retrieve recordings")
		return
	}

	coreUtil.RespondSuccess(c, recordings)
}


// GET /api/recordings/{id}
// @Summary Get a specific recording by ID for the logged-in user
// @Description Retrieve a recording by its ID for the logged-in user
// @Tags recordings
// @Accept json
// @Produce json
// @Param id path string true "Recording ID"
// @Success 200 {object} model.Recording "The requested recording"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 404 {object} map[string]string "Recording not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/recordings/{id} [get]
func (h *RecordingHandler) GetRecordingByID(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := coreUtil.GetUserIDFromContext(c)
	if err != nil {
		coreUtil.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		coreUtil.RespondUnauthorized(c, "Invalid user ID")
		return
	}

	recordingID := c.Param("id")

	recording, err := h.service.GetRecordingByID(ctx, userID, recordingID)
	if err != nil {
		if err.Error() == "recording not found" {
			coreUtil.RespondNotFound(c, "Recording not found")
		} else {
			coreUtil.RespondInternalError(c, "Failed to retrieve recording")
		}
		return
	}

	coreUtil.RespondSuccess(c, recording)
}