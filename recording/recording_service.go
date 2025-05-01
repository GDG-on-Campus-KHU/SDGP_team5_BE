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
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

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


type RecordingService struct{}

func (s *RecordingService) GetRecordingsByUserID(ctx context.Context, userID int) ([]model.Recording, error) {
	var recordings []model.Recording

	filter := bson.M{"user_id": userID}
	cursor, err := dbConfig.RecordingCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var recording model.Recording
		if err := cursor.Decode(&recording); err != nil {
			return nil, err
		}
		recordings = append(recordings, recording)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return recordings, nil
}


func (s *RecordingService) GetRecordingByID(ctx context.Context, userID int, objectID string) (*model.Recording, error) {
	objID, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		return nil, errors.New("invalid recording id format")
	}

	filter := bson.M{
		"_id":     objID,
		"user_id": userID,
	}

	var recording model.Recording
	err = dbConfig.RecordingCollection.FindOne(ctx, filter).Decode(&recording)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("recording not found")
		}
		return nil, err
	}

	return &recording, nil
}