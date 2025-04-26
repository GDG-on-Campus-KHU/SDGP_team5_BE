// recording/recording_service.go

package recording

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/recording/util"
)

// recording data
func NewRecording(userID int, countryCode string, recordingURL string) *model.Recording {
	recordingID := generateRecordingID(userID)

	return &model.Recording{
		RecordingID:   recordingID,
		UserID:        userID,
		RecordingURL:  recordingURL,
		RecordingText: "",				// 'NULL' by default until SpeechToText result is available
		CountryCode:   countryCode,
		CreatedAt:     time.Now(),
	}
}

// userID와 현재 시간으로 'RecordingID' 생성 (unique)
func generateRecordingID(userID int) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("user_%s_%s", strconv.Itoa(userID), timestamp)
}


func SaveRecording(ctx context.Context, userID int, countryCode string, file *multipart.FileHeader) (*model.Recording, error) {
	// upload file to GCS
	recordingURL, err := util.UploadFileToGCS(file, fmt.Sprintf("%d", userID))
	if err != nil {
		log.Printf("GCS Upload Error: %v", err)
		return nil, fmt.Errorf("failed to upload file to GCS: %v", err)
	}

	// create recording metadata
	recording := NewRecording(userID, countryCode, recordingURL)

	// integrate with MongoDB collection
	collection := dbConfig.RecordingCollection

	// insert recording metadata into MongoDB
	_, err = collection.InsertOne(ctx, recording)
	if err != nil {
		log.Printf("MongoDB Insert Error: %v", err)
		return nil, fmt.Errorf("failed to insert recording into MongoDB: %v", err)
	}

	return recording, nil
}