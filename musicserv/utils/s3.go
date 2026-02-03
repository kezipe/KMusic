package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitS3() error {
	endpoint := os.Getenv("S3_ENDPOINT") // e.g., "s3.amazonaws.com"
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	useSSL := true

	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create minio client: %v", err)
	}

	minioClient = client
	return nil
}

// UploadFile uploads a file to S3 and returns the object name (key)
func UploadFile(bucketName string, objectName string, reader io.Reader, contentType string) error {
	ctx := context.Background()

	// Check if bucket exists
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	// Upload the file
	_, err = minioClient.PutObject(ctx, bucketName, objectName, reader, -1,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// GetPresignedURL generates a presigned URL for an object with 15-minute expiration
func GetPresignedURL(bucketName, objectName string) (string, error) {
	ctx := context.Background()

	// Generate presigned URL valid for 15 minutes
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Minute*15, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.String(), nil
}

// DeleteObject deletes an object from S3
func DeleteObject(bucketName, objectName string) error {
	ctx := context.Background()

	err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %v", err)
	}

	return nil
}
