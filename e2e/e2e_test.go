//go:build e2e

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/saltpay/go-kafka-driver"
)

type validTest struct {
	F1 string `json:"f1"`
}

func TestValidator(t *testing.T) {
	var (
		ctx           = context.Background()
		kafkaEndpoint = os.Getenv("CONSUMER_KAFKA_ENDPOINT")
		kafkaUsername = os.Getenv("CONSUMER_KAFKA_USERNAME")
		kafkaPassword = os.Getenv("CONSUMER_KAFKA_PASSWORD")
		kafkaIngress  = os.Getenv("CONSUMER_KAFKA_TOPIC")
		kafkaEgress   = os.Getenv("PRODUCER_KAFKA_TOPIC")
		payload       = validTest{F1: "field"}
	)

	// producer
	producer, err := kafka.NewProducer(ctx, kafka.ProducerConfig{
		Addr:     strings.Split(kafkaEndpoint, ","),
		Topic:    kafkaIngress,
		Username: kafkaUsername,
		Password: kafkaPassword,
	})
	if err != nil {
		t.Fatal(err)
	}

	consumer, err := kafka.NewConsumer(ctx, kafka.ConsumerConfig{
		Brokers:  strings.Split(kafkaEndpoint, ","),
		GroupID:  fmt.Sprintf("%s-groupID", kafkaEgress),
		Topic:    kafkaEgress,
		Username: kafkaUsername,
		Password: kafkaPassword,
	})
	if err != nil {
		t.Fatal(err)
	}

	is := is.New(t)
	t.Run("An incoming transaction in the kafka topic should generate a transaction in the DB", func(t *testing.T) {
		// given a processor
		done := make(chan bool)
		var msgJSON []byte
		processor := func(ctx context.Context, msg kafka.Message) error {
			msgJSON = msg.Value
			done <- true
			return nil
		}
		pauseStrategy := func(ctx context.Context) bool {
			return true
		}

		// when a message is written
		payloadJSON, err := json.Marshal(payload)
		is.NoErr(err)
		err = producer.WriteMessage(ctx, kafka.Message{
			Value: payloadJSON,
		})
		is.NoErr(err)

		go consumer.Listen(ctx, processor, kafka.AlwaysCommitWithoutError, pauseStrategy)

		select {
		case <-done:
			is.Equal(msgJSON, payloadJSON)
		case <-time.After(5 * time.Second):
			t.Fatal("time out while waiting for test to finish")
		}
	})
}
