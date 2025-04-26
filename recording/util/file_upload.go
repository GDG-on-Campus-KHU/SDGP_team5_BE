// recording/util/file_upload.go

package util

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/storage"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
)

// GCS bucket에 음성 녹음 파일 업로드하고 file URL을 return
func UploadFileToGCS(file *multipart.FileHeader, userID string) (string, error) {
	// GCS client
	ctx := context.Background()
	client := dbConfig.GCSClient

	// GCS bucket
	bucketName := "resq-upload-bucket"
	bucket := client.Bucket(bucketName)

	// unique filename generation rule
	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	timestamp := time.Now().Format("060102150405")	// yyMMddhhmmss
	fileName := fmt.Sprintf("user_%s_%s%s", userID, timestamp, ext)

	srcFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer srcFile.Close()

	// file upload
	dst := bucket.Object(fileName).NewWriter(ctx)
	dst.ContentType = file.Header.Get("Content-Type")
	if _, err := io.Copy(dst, srcFile); err != nil {
		return "", fmt.Errorf("failed to copy file to GCS: %v", err)
	}
	
	if err := dst.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %v", err)
	}

	// public read access
	object := bucket.Object(fileName)
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set ACL for public read: %v", err)
	}

	// return GCS file URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return url, nil
}