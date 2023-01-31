//go:build unit

package kafkalistener_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/saltpay/go-kafka-driver"

	kafkalistener "github.com/saltpay/transaction-card-validator/internal/adapters/kafka"
	"github.com/saltpay/transaction-card-validator/internal/application/ports"
	"github.com/saltpay/transaction-card-validator/internal/application/ports/mocks"
)

func TestListeners(t *testing.T) {
	is := is.New(t)
	t.Run("valid messages can skip validation and be written by producer", func(t *testing.T) {
		producerDone := make(chan bool)

		mockProducer := mocks.ProducerMock{WriteMessageFunc: func(ctx context.Context, msg kafka.Message) error {
			producerDone <- true
			return nil
		}}
		mockConsumer := mocks.ConsumerMock{ListenFunc: func(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy) {
			err := processor(ctx, kafka.Message{})
			is.NoErr(err)
		}}
		mockS3Client := mocks.QuarantineHandlerMock{UploadObjectFunc: func(topicName string, data []byte) error {
			return nil
		}}
		mockValidator := &mocks.ValidatorMock{ValidateFunc: func(msgDecoded interface{}) (bool, error) {
			return true, nil
		}}
		mockRegistry := &mocks.SchemaHandlerMock{DecodeFunc: func(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
			return nil, true
		}}

		listener := kafkalistener.Listener{
			Producer:            &mockProducer,
			Consumer:            &mockConsumer,
			S3Client:            &mockS3Client,
			BusinessValidator:   mockValidator,
			SchemaHandlerClient: mockRegistry,
			MetricsClient:       &ports.DummyMetricsClient{},
			SchemaKey:           "key",
			TopicName:           "topic",
			SkipValidation:      true,
			WorkerPoolSize:      1,
		}

		go listener.Listen(context.TODO())

		select {
		case <-producerDone:
			is.Equal(len(mockProducer.WriteMessageCalls()), 1)
			is.Equal(len(mockValidator.ValidateCalls()), 0)
			is.Equal(len(mockS3Client.UploadObjectCalls()), 0)
			is.Equal(len(mockRegistry.DecodeCalls()), 1)
		case <-time.After(50 * time.Millisecond):
			t.Fatal("time out while waiting for test to finish")
		}
	})
	t.Run("valid messages should be written by the producer", func(t *testing.T) {
		producerDone := make(chan bool)

		mockProducer := mocks.ProducerMock{WriteMessageFunc: func(ctx context.Context, msg kafka.Message) error {
			producerDone <- true
			return nil
		}}
		mockConsumer := mocks.ConsumerMock{ListenFunc: func(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy) {
			err := processor(ctx, kafka.Message{})
			is.NoErr(err)
		}}
		mockS3Client := mocks.QuarantineHandlerMock{UploadObjectFunc: func(topicName string, data []byte) error {
			return nil
		}}
		mockValidator := &mocks.ValidatorMock{ValidateFunc: func(msgDecoded interface{}) (bool, error) {
			return true, nil
		}}
		mockRegistry := &mocks.SchemaHandlerMock{DecodeFunc: func(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
			return nil, true
		}}

		listener := kafkalistener.Listener{
			Producer:            &mockProducer,
			Consumer:            &mockConsumer,
			S3Client:            &mockS3Client,
			BusinessValidator:   mockValidator,
			SchemaHandlerClient: mockRegistry,
			MetricsClient:       &ports.DummyMetricsClient{},
			SchemaKey:           "key",
			TopicName:           "topic",
			WorkerPoolSize:      1,
		}

		go listener.Listen(context.TODO())

		select {
		case <-producerDone:
			is.Equal(len(mockProducer.WriteMessageCalls()), 1)
			is.Equal(len(mockValidator.ValidateCalls()), 1)
			is.Equal(len(mockS3Client.UploadObjectCalls()), 0)
			is.Equal(len(mockRegistry.DecodeCalls()), 1)
		case <-time.After(50 * time.Millisecond):
			t.Fatal("time out while waiting for test to finish")
		}
	})
	t.Run("if validation fails, message should be quarantined and not written back", func(t *testing.T) {
		consumerDone := make(chan bool)
		validatorDone := make(chan bool)
		s3ClientDone := make(chan bool)

		mockProducer := mocks.ProducerMock{WriteMessageFunc: func(ctx context.Context, msg kafka.Message) error {
			return nil
		}}
		mockConsumer := mocks.ConsumerMock{ListenFunc: func(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy) {
			err := processor(ctx, kafka.Message{})
			is.NoErr(err)
			consumerDone <- true
		}}
		mockS3Client := mocks.QuarantineHandlerMock{UploadObjectFunc: func(topicName string, data []byte) error {
			s3ClientDone <- true
			return nil
		}}
		mockValidator := &mocks.ValidatorMock{ValidateFunc: func(msgDecoded interface{}) (bool, error) {
			validatorDone <- true
			return false, errors.New("validation failed")
		}}
		mockRegistry := &mocks.SchemaHandlerMock{DecodeFunc: func(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
			return nil, true
		}}

		listener := kafkalistener.Listener{
			Producer:            &mockProducer,
			Consumer:            &mockConsumer,
			S3Client:            &mockS3Client,
			BusinessValidator:   mockValidator,
			SchemaHandlerClient: mockRegistry,
			MetricsClient:       &ports.DummyMetricsClient{},
			SchemaKey:           "key",
			TopicName:           "topic",
			WorkerPoolSize:      1,
		}

		go listener.Listen(context.TODO())

		select {
		case <-consumerDone:
		case <-validatorDone:
		case <-s3ClientDone:
			is.Equal(len(mockProducer.WriteMessageCalls()), 0)
			is.Equal(len(mockValidator.ValidateCalls()), 1)
			is.Equal(len(mockS3Client.UploadObjectCalls()), 1)
			is.Equal(len(mockRegistry.DecodeCalls()), 1)
		case <-time.After(50 * time.Millisecond):
			t.Fatal("time out while waiting for test to finish")
		}
	})
}
