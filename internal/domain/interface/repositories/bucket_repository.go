package repositories

import (
	"context"
	"mime/multipart"
)

type BucketRepository interface {
	CreateFile(ctx context.Context, path string, filename string, file multipart.File) (string, error)
}
