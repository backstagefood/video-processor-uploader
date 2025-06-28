package repositories

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/backstagefood/video-processor-uploader/internal/domain/interface/repositories"
	"github.com/backstagefood/video-processor-uploader/pkg/adapter/bucketconfig"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type bucketRepository struct {
	s3Conn     *s3.S3
	bucketName string
}

func NewBucketRepository(s3Conn *bucketconfig.ApplicationS3Bucket) repositories.BucketRepository {
	return &bucketRepository{
		s3Conn:     s3Conn.Client(),
		bucketName: s3Conn.BucketName(),
	}
}

func (v *bucketRepository) CreateFile(ctx context.Context, path string, filename string, file multipart.File) (string, error) {
	// Ensure clean path construction
	key := filepath.Join(path, filename)
	// For S3, we need forward slashes regardless of OS
	key = filepath.ToSlash(key)
	// Remove any leading/trailing slashes that might cause issues
	key = strings.Trim(key, "/")

	slog.Info("bucketRepository", "key", key)

	// CreateFile directly from the file reader without loading into memory
	_, err := v.s3Conn.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(v.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		slog.Error("não foi possível subir o arquivo", "error", err)
		return "", fmt.Errorf("não foi possível subir o arquivo: %w", err)
	}

	return key, nil
}
