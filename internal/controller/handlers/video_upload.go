package handlers

import (
	"fmt"
	"github.com/backstagefood/video-processor-uploader/internal/domain"
	"github.com/backstagefood/video-processor-uploader/internal/domain/interface/adapters"
	portServices "github.com/backstagefood/video-processor-uploader/internal/domain/interface/services"
	"github.com/backstagefood/video-processor-uploader/internal/repositories"
	"github.com/backstagefood/video-processor-uploader/internal/usecase"
	"github.com/backstagefood/video-processor-uploader/pkg/adapter/bucketconfig"
	"github.com/backstagefood/video-processor-uploader/utils"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type VideoHandler struct {
	filesService portServices.FilesService
}

func NewVideoHandler(messageProducer adapters.MessageProducer, s3Client *bucketconfig.ApplicationS3Bucket) *VideoHandler {
	bucketRepository := repositories.NewBucketRepository(s3Client)
	return &VideoHandler{
		filesService: usecase.NewFileService(messageProducer, bucketRepository),
	}
}

// @BasePath /v1/upload
// PingExample godoc
// @Summary Upload video file
// @Schemes
// @Description Upload video file in the following formats: mp4, avi, mov, mkv
// @Tags upload
// @Accept multipart/form-data
// @Param video formData file true "Video file to upload"
// @Produce application/json
// @Success 200 {object} object{success=boolean, message=string} "success response"
// @Failure 500 {object} object{error=string} "generic error response"
// @Router /v1/upload [post]
func (v *VideoHandler) HandleVideoUpload(c *gin.Context) {
	userEmail := c.MustGet("user_email").(string)
	slog.Info("obtem userEmail em handleVideoUpload", "userEmail", userEmail)

	file, header, err := c.Request.FormFile("video")
	if err != nil {
		var processingResult = getErrorProcessingResult("Erro ao receber arquivo: " + err.Error())
		c.JSON(processingResult.Code, processingResult)
		return
	}
	defer file.Close()

	if !utils.IsValidVideoFile(header.Filename) {
		var processingResult = getErrorProcessingResult("Formato de arquivo não suportado. Use: mp4, avi, mov, mkv")
		c.JSON(processingResult.Code, processingResult)
		return
	}

	// gravar no bucket s3
	_, err = v.filesService.CreateFile(c, userEmail, header.Filename, file)
	if err != nil {
		var processingResult = getErrorProcessingResult(err.Error())
		c.JSON(processingResult.Code, processingResult)
		return
	}
	c.JSON(200, domain.ProcessingResult{
		Success: true,
		Code:    200,
		Message: fmt.Sprintf("estamos processando seu vídeo, aguarde um momento"),
	})
}

func getErrorProcessingResult(message string) domain.ProcessingResult {
	var processingResult = domain.ProcessingResult{
		Code:    400,
		Success: false,
		Message: message,
	}
	return processingResult
}
