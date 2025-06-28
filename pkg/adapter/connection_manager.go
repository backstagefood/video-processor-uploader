package adapter

import (
	"github.com/backstagefood/video-processor-uploader/internal/domain/interface/adapters"
	"github.com/backstagefood/video-processor-uploader/pkg/adapter/bucketconfig"
	"github.com/backstagefood/video-processor-uploader/pkg/adapter/kafka"
	"log/slog"
	"os"
)

type ConnectionManager interface {
	GetBucketConn() *bucketconfig.ApplicationS3Bucket
	GetMessageProducer() adapters.MessageProducer
}

type connectionManagerImpl struct {
	bucketConn      *bucketconfig.ApplicationS3Bucket
	messageProducer adapters.MessageProducer
}

func NewConnectionManager() ConnectionManager {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	producer, err := kafka.NewProducer(broker, topic)
	if err != nil {
		slog.Error("não foi possível criar o produtor do topico kafka", "error", err)
	}
	return &connectionManagerImpl{
		bucketConn:      bucketconfig.NewBucketConnection(),
		messageProducer: producer,
	}
}

func (c *connectionManagerImpl) GetBucketConn() *bucketconfig.ApplicationS3Bucket {
	return c.bucketConn
}

func (c *connectionManagerImpl) GetMessageProducer() adapters.MessageProducer {
	return c.messageProducer
}
