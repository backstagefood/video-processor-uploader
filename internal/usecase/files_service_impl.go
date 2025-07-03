package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/backstagefood/video-processor-uploader/internal/domain"
	"github.com/backstagefood/video-processor-uploader/internal/domain/interface/adapters"
	portRepositories "github.com/backstagefood/video-processor-uploader/internal/domain/interface/repositories"
	portServices "github.com/backstagefood/video-processor-uploader/internal/domain/interface/services"
	"github.com/backstagefood/video-processor-uploader/utils"
	"log/slog"
	"mime/multipart"
	"path/filepath"
)

type fileService struct {
	messageProducer  adapters.MessageProducer
	bucketRepository portRepositories.BucketRepository
}

func NewFileService(
	messageProducer adapters.MessageProducer,
	bucketRepository portRepositories.BucketRepository) portServices.FilesService {
	return &fileService{
		messageProducer:  messageProducer,
		bucketRepository: bucketRepository,
	}
}

func (f *fileService) CreateFile(ctx context.Context, userEmail string, fileName string, file multipart.File) (string, error) {
	slog.Info("fileService - create file", "userEmail", userEmail, "fileName", fileName)
	// junta nome do usuario com caminho
	path := filepath.Join(utils.SanitizeEmailForPath(userEmail), "video_files")

	// grava no bucket
	//timestamp := time.Now().Format("20060102_150405")
	filePrefix := utils.GenerateUniqueKey()
	fileNameWithTimestamp := fmt.Sprintf("%s_%s", filePrefix, fileName)
	fileFullPath, err := f.bucketRepository.CreateFile(ctx, path, fileNameWithTimestamp, file)
	if err != nil {
		slog.Error("não foi possível gravar o video no bucket", "error", err)
		return "", err
	}

	fileSize, err := utils.GetFileSize(file)
	if err != nil {
		slog.Error("não foi possível obter o tamanho do vídeo", "error", err)
		return "", err
	}

	// publica no topico kafka
	eventBytes, _ := json.Marshal(&domain.FilePayload{UserName: userEmail, FilePath: fileFullPath, FileSize: fileSize})
	err = f.messageProducer.ProduceMessage(ctx, utils.GetEnvVarOrDefault("KAFKA_TOPIC", ""), eventBytes)
	if err != nil {
		slog.Error("não foi enviar o vídeo pelo tópico kafka", "error", err)
		return "", err
	}
	return fileFullPath, nil
}
