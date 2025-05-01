// recording/recording_service.go

package recording

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"time"
	"io"
	"os"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/recording/util"
	coreUtil "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)

// recording data
func NewRecording(userID int, recordingURL string) *model.Recording {
	recordingID := generateRecordingID(userID)

	return &model.Recording{
		RecordingID:   recordingID,
		UserID:        userID,
		RecordingURL:  recordingURL,
		CreatedAt:     time.Now(),
	}
}

// userID와 현재 시간으로 'RecordingID' 생성 (unique)
func generateRecordingID(userID int) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("user_%s_%s", strconv.Itoa(userID), timestamp)
}


func ShortRecordingService(ctx context.Context, userID int, file *multipart.FileHeader) (string, string, error) {
	user, err := coreUtil.GetUserByUserID(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %v", err)
	}

	openedFile, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open file: %v", err)
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}

	// temporary file
	tempFilename := fmt.Sprintf("tmp/%d_%s", time.Now().Unix(), file.Filename)
	err = os.WriteFile(tempFilename, fileBytes, 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write temp file: %v", err)
	}
	defer os.Remove(tempFilename)  // clean up temporary file

	languageCode := coreUtil.AppLangToLangCode(user.AppLang)

	// request synchonously Speech-to-Text
	sttResult, err := util.SynchronousSpeechToText(tempFilename, languageCode)
	if err != nil {
		return "", "", fmt.Errorf("speech to text failed: %v", err)
	}

	log.Printf("STT Result: %v", sttResult) 

	// file upload to GCS
	recordingURL, err := util.UploadFileToGCS(file, fmt.Sprintf("%d", userID))
	if err != nil {
		log.Printf("file upload failed after STT: %v", err)
		return "", "", fmt.Errorf("file upload failed: %v", err)
	}

	// save to database
	collection := dbConfig.RecordingCollection
	_, err = collection.InsertOne(ctx, map[string]interface{}{
		"user_id":        userID,
		"recording_url":  recordingURL,
		"recording_text": sttResult,
		"app_lang":       user.AppLang,
		"created_at":     time.Now(),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to insert recording into DB: %v", err)
	}

	return recordingURL, sttResult, nil
}