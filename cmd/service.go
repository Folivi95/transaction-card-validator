package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/saltpay/go-kafka-driver"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-card-validator/internal/adapters/http"
	kafkalistener "github.com/saltpay/transaction-card-validator/internal/adapters/kafka"
	"github.com/saltpay/transaction-card-validator/internal/adapters/s3"
	"github.com/saltpay/transaction-card-validator/internal/adapters/schemaregistry"
	"github.com/saltpay/transaction-card-validator/internal/adapters/testhelpers"
	pgatewayvalidation "github.com/saltpay/transaction-card-validator/internal/application/pgateway"
	"github.com/saltpay/transaction-card-validator/internal/application/ports"
	way4validation "github.com/saltpay/transaction-card-validator/internal/application/way4"
)

type Service struct {
	Listener     kafkalistener.Listener
	ServerConfig http.ServerConfig
}

func newService(ctx context.Context) (*Service, error) {
	httpConfig := newServerConfig()

	appConfig, err := loadAppConfig()
	if err != nil {
		return &Service{}, fmt.Errorf("failed to load configs")
	}

	// load and initialize log level
	InitLog(appConfig)

	// Create Kafka Clients
	kafkaConsumerConfig := kafka.ConsumerConfig{
		Brokers:  appConfig.KafkaConsumerConfig.KafkaEndpoint,
		GroupID:  fmt.Sprintf("%s-groupID", appConfig.KafkaConsumerConfig.KafkaTopic),
		Topic:    appConfig.KafkaConsumerConfig.KafkaTopic,
		Username: appConfig.KafkaConsumerConfig.KafkaUsername,
		Password: appConfig.KafkaConsumerConfig.KafkaPassword,
	}
	kafkaConsumer, err := kafka.NewConsumer(ctx, kafkaConsumerConfig)
	if err != nil {
		return &Service{}, err
	}

	kafkaProducerConfig := kafka.ProducerConfig{
		Addr:     appConfig.KafkaProducerConfig.KafkaEndpoint,
		Topic:    appConfig.KafkaProducerConfig.KafkaTopic,
		Username: appConfig.KafkaProducerConfig.KafkaUsername,
		Password: appConfig.KafkaProducerConfig.KafkaPassword,
	}
	kafkaProducer, err := kafka.NewProducer(ctx, kafkaProducerConfig)
	if err != nil {
		return &Service{}, err
	}

	// Create S3 Client
	s3Config := aws.Config{
		Endpoint:         aws.String(appConfig.S3Config.S3Endpoint),
		Region:           aws.String(appConfig.S3Config.S3Region),
		DisableSSL:       aws.Bool(appConfig.S3Config.S3DisableSSL),
		S3ForcePathStyle: aws.Bool(appConfig.S3Config.S3ForcePathStyle),
	}
	s3Client, err := s3.NewS3Client(&s3Config, appConfig.S3Config.S3Bucket)
	if err != nil {
		return &Service{}, err
	}

	// Get Schema Client - Local or Registry
	var srClient ports.SchemaHandler
	if appConfig.SchemaRegistryConfig.RegistryEndpoint != "" {
		srClient = schemaregistry.NewSchemaRegistryClient(appConfig.SchemaRegistryConfig.RegistryEndpoint, appConfig.SchemaRegistryConfig.RegistryRefreshTimeSeconds)
	} else {
		srClient = testhelpers.NewLocalModel()
	}

	var validator ports.Validator
	businessValidationName := os.Getenv("BUSINESS_VALIDATION")
	switch businessValidationName {
	case "Way4MaskedCardNumber":
		cardHoldKeys := strings.Split(os.Getenv("W4_MASKED_CHD_KEYS"), ",")
		maskRegexPattern := os.Getenv("W4_MASKED_REGEX_PATTERN")
		validator, err = way4validation.NewWay4MaskedCardNumber(cardHoldKeys, maskRegexPattern)
		if err != nil {
			return &Service{}, err
		}
	case "PGatewayValidation":
		validator, err = pgatewayvalidation.NewValidateNoCHD()
		if err != nil {
			return &Service{}, err
		}
	default:
		validator = &ports.DummyValidator{}
	}

	// Skip validation if SKIP_VALIDATION is set
	skipValidation, err := strconv.ParseBool(os.Getenv("SKIP_VALIDATION"))
	if err != nil {
		zapctx.Warn(ctx, "[processorHandler] Failed to get SKIP_VALIDATION parameter. Default to FALSE.", zap.Error(err))
		skipValidation = false
	}

	// Stop processing if PAUSE_PROCESSING is set
	pauseProcessing, err := strconv.ParseBool(os.Getenv("PAUSE_PROCESSING"))
	if err != nil {
		zapctx.Warn(ctx, "[processorHandler] Failed to get PAUSE_PROCESSING flag. Default to FALSE.", zap.Error(err))
		pauseProcessing = false
	}

	pss, ok := os.LookupEnv("WORKER_POOL_SIZE")
	if !ok {
		pss = "150"
	}
	workerPoolSize, err := strconv.Atoi(pss)
	if err != nil {
		zapctx.Warn(ctx, "[startWorkerGroup] Failed to set worker pool size from env. Default to 150.", zap.Error(err))
	}

	// Patching context with srClient and key, producer and s3Client
	validatorListener := kafkalistener.Listener{
		SchemaKey:           appConfig.SchemaRegistryConfig.RegistrySubjectName,
		Producer:            kafkaProducer,
		Consumer:            kafkaConsumer,
		S3Client:            s3Client,
		BusinessValidator:   validator,
		SchemaHandlerClient: srClient,
		TopicName:           appConfig.KafkaConsumerConfig.KafkaTopic,
		AuditLogFields:      appConfig.LogConfig.AuditLogFields,
		MetricsClient:       newMetricsClient(),
		SkipValidation:      skipValidation,
		PauseProcessing:     pauseProcessing,
		WorkerPoolSize:      workerPoolSize,
	}

	return &Service{
		Listener:     validatorListener,
		ServerConfig: httpConfig,
	}, nil
}
