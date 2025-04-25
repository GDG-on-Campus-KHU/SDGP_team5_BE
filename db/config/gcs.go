// db/config/gcs.go

package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var GCSClient *storage.Client

func InitGCS() {
	credentialsPath := os.Getenv("GOOGLE_CREDENTIALS")

	if credentialsPath == "" {
		log.Fatal("ERROR: GOOGLE_CREDENTIALS is NOT set.")
	}

	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("ERROR: failed to create GCS client: %v", err)
	}

	GCSClient = client

	fmt.Println("Successfully initialized GCS client.")
}
