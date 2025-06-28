package services

import (
	"context"
	"mime/multipart"
)

type FilesService interface {
	CreateFile(ctx context.Context, userEmail string, filename string, file multipart.File) (string, error)
}
