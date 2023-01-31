package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/saltpay/transaction-card-validator/internal/adapters/http"
)

const (
	CONSUMER = "CONSUMER"
	PRODUCER = "PRODUCER"
)

type AppConfig struct {
	KafkaProducerConfig  KafkaConfig
	KafkaConsumerConfig  KafkaConfig
	SchemaRegistryConfig SchemaRegistryConfig
	S3Config             S3Config
	LogConfig            LogConfig
}

type LogConfig struct {
	LogLevel       string   `split_words:"true"`
	AuditLogFields []string `split_words:"true"`
}

type KafkaConfig struct {
	KafkaEndpoint []string `split_words:"true"`
	KafkaUsername string   `split_words:"true"`
	KafkaPassword string   `split_words:"true"`
	KafkaTopic    string   `split_words:"true"`
}

type S3Config struct {
	S3Bucket         string `split_words:"true"`
	S3Endpoint       string `split_words:"true"`
	S3Region         string `split_words:"true"`
	S3DisableSSL     bool   `split_words:"true"`
	S3ForcePathStyle bool   `split_words:"true"`
}

type SchemaRegistryConfig struct {
	RegistryEndpoint           string `split_words:"true"`
	RegistrySubjectName        string `split_words:"true"`
	RegistryRefreshTimeSeconds int    `split_words:"true"`
}

func loadAppConfig() (AppConfig, error) {
	kafkaConsumerConfig, err := loadKafkaConfig(CONSUMER)
	if err != nil {
		return AppConfig{}, err
	}

	kafkaProducerConfig, err := loadKafkaConfig(PRODUCER)
	if err != nil {
		return AppConfig{}, err
	}

	s3Config, err := loadS3Config()
	if err != nil {
		return AppConfig{}, err
	}

	logConfig, err := loadLogConfig()
	if err != nil {
		return AppConfig{}, err
	}

	schemaRegistryConfig, err := loadSchemaRegistryConfig()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		KafkaConsumerConfig:  kafkaConsumerConfig,
		KafkaProducerConfig:  kafkaProducerConfig,
		SchemaRegistryConfig: schemaRegistryConfig,
		S3Config:             s3Config,
		LogConfig:            logConfig,
	}, nil
}

// LoadKafkaConfig loads the app config from environment variables.
func loadKafkaConfig(prefix string) (KafkaConfig, error) {
	var config KafkaConfig
	err := envconfig.Process(prefix, &config)
	if err != nil {
		return KafkaConfig{}, err
	}
	return config, nil
}

func loadS3Config() (S3Config, error) {
	var config S3Config
	err := envconfig.Process("", &config)
	if err != nil {
		return S3Config{}, err
	}
	return config, nil
}

func loadSchemaRegistryConfig() (SchemaRegistryConfig, error) {
	var config SchemaRegistryConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return SchemaRegistryConfig{}, err
	}

	return config, nil
}

func loadLogConfig() (LogConfig, error) {
	var logLevel LogConfig
	err := envconfig.Process("", &logLevel)
	if err != nil {
		return LogConfig{}, err
	}

	return logLevel, nil
}

func newServerConfig() http.ServerConfig {
	return http.ServerConfig{
		Port:             "8080",
		HTTPReadTimeout:  2 * time.Second,
		HTTPWriteTimeout: 2 * time.Second,
	}
}
