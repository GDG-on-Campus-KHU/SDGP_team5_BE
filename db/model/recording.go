// db/model/recording.go

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recording struct {
	ID            primitive.ObjectID	`bson:"_id,omitempty" json:"id,omitempty"`                 // MongoDB ObjectId
	RecordingID   string				`bson:"recording_id" json:"recording_id"`                  // auto-increment ID
	UserID        int                	`bson:"user_id" json:"user_id"`                            // ref. struct 'User'
	RecordingURL  string             	`bson:"recording_url" json:"recording_url"`                // 음성 파일 GCS URL
	RecordingText string             	`bson:"recording_text,omitempty" json:"recording_text"`    // SpeechToText 결과물 (NULL by default)
	AppLang       string             	`bson:"app_lang" json:"app_lang"`  
	CreatedAt     time.Time			 	`bson:"created_at,omitempty" json:"created_at,omitempty"`	// timestamp
}