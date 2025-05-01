// recording/util/synchronous_stt.go

package util

import (
	"context"
	"fmt"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/api/option"
)

const (
	GoogleSTTAPIEndpoint = "https://speech.googleapis.com/v1/speech:recognize"
)


func SynchronousSpeechToText(filePath string, languageCode string) (string, error) {
	
	// GOOGLE_CREDENTIALS
	credentialsPath := os.Getenv("GOOGLE_CREDENTIALS")
	if credentialsPath == "" {
		return "", fmt.Errorf("GOOGLE_CREDENTIALS not set in environment")
	}

	// GCS client
	ctx := context.Background()
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return "", fmt.Errorf("Failed to create Speech client: %v", err)
	}
	defer client.Close()

	data, err := os.ReadFile(filePath)

	if err != nil {
		return "", fmt.Errorf("Failed to read audio file: %v", err)
	}

	// RecognizeRequest 생성
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			// SampleRateHertz: 16000,
			SampleRateHertz: 44100,  // audio sample rate
			LanguageCode: languageCode,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	// request
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("Failed to recognize: %v", err)
	}

	// result
	var transcription string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			transcription += alt.Transcript
		}
	}

	if transcription == "" {
		return "", fmt.Errorf("No transcription found")
	}

	return transcription, nil
}