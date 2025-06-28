package adapters

import (
	"context"
)

type MessageProducer interface {
	ProduceMessage(ctx context.Context, key string, value []byte) error
}
