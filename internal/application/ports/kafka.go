//go:generate moq -out mocks/kafka_consumer_moq.go -pkg=mocks . Consumer
//go:generate moq -out mocks/kafka_producer_moq.go -pkg=mocks . Producer

package ports

import (
	"context"

	"github.com/saltpay/go-kafka-driver"
)

type Consumer interface {
	Listen(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy)
}

type Producer interface {
	WriteMessage(ctx context.Context, msg kafka.Message) error
}
