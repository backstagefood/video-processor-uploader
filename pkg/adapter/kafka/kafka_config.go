package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/backstagefood/video-processor-uploader/internal/domain/interface/adapters"
	"log/slog"
)

func newKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	return config
}

type Producer struct {
	Producer sarama.SyncProducer
	Topic    string
}

func NewProducer(broker string, topic string) (adapters.MessageProducer, error) {
	producer, err := sarama.NewSyncProducer([]string{broker}, newKafkaConfig())
	if err != nil {
		slog.Error("error with kafka sync producer", slog.String("error", err.Error()))
		return nil, err
	}
	return &Producer{Producer: producer, Topic: topic}, nil
}

func (kp *Producer) ProduceMessage(ctx context.Context, key string, value []byte) error {
	valueString := string(value)
	msg := &sarama.ProducerMessage{
		Topic: kp.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(valueString),
	}
	slog.InfoContext(
		ctx,
		"sending message",
		slog.Group(
			"message",
			slog.String("key", key),
			slog.String("value", valueString),
		),
	)
	_, _, err := kp.Producer.SendMessage(msg)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"error producing message",
			err.Error(),
			slog.Group(
				"message",
				slog.String("key", key),
				slog.String("value", valueString),
			),
		)
	}
	return err
}
