package kafkalistener

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/saltpay/go-kafka-driver"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-card-validator/internal/application/ports"
)

type Listener struct {
	Producer            ports.Producer
	Consumer            ports.Consumer
	S3Client            ports.QuarantineHandler
	BusinessValidator   ports.Validator
	SchemaHandlerClient ports.SchemaHandler
	MetricsClient       ports.MetricsClient
	SchemaKey           string
	TopicName           string
	AuditLogFields      []string
	SkipValidation      bool
	PauseProcessing     bool
	MessageChannel      chan kafka.Message
	WorkerPoolSize      int
}

func (l *Listener) Listen(ctx context.Context) {
	l.startWorkerGroup(ctx)
	l.Consumer.Listen(ctx, l.onMessage, kafka.AlwaysCommitWithoutError, l.pauseProcessing)
}

// todo: this could be a feature flag that we can change while the app is live
func (l *Listener) pauseProcessing(ctx context.Context) bool {
	return l.PauseProcessing
}

func (l *Listener) quarantineHandler(ctx context.Context, message []byte) {
	dynamic := make(map[string]interface{})
	err := json.Unmarshal(message, &dynamic)
	if err != nil {
		zapctx.Warn(ctx, "[Audit] Failed to unmarshal incoming message. Sending raw JSON to quarantine", zap.Error(err))
		if err := l.S3Client.UploadObject(l.TopicName, message); err != nil {
			zapctx.Warn(ctx, "[Audit] Failed to write message to quarantine.", zap.Error(err))
		}
		return
	}

	// Get card_hold keys
	cardHoldKeys := strings.Split(os.Getenv("CARD_HOLD_KEYS"), ",")

	for _, cardHoldKey := range cardHoldKeys {
		// Remove card hold data
		if dynamic["before"] != nil {
			delete(dynamic["before"].(map[string]interface{}), cardHoldKey)
		}

		if dynamic["after"] != nil {
			delete(dynamic["after"].(map[string]interface{}), cardHoldKey)
		}
	}

	msgBytes, _ := json.Marshal(dynamic)

	if err := l.S3Client.UploadObject(l.TopicName, msgBytes); err != nil {
		zapctx.Warn(ctx, "[Audit] Failed to write message to quarantine.", zap.Error(err))
	}
}

func (l *Listener) startWorkerGroup(ctx context.Context) {
	l.MessageChannel = make(chan kafka.Message)

	for w := 0; w < l.WorkerPoolSize; w++ {
		go func() {
			for {
				err := l.processMessage(ctx, <-l.MessageChannel)
				if err != nil {
					zapctx.Error(ctx, "[Processor] Unable to process message.", zap.Error(err))
				}
			}
		}()
	}
}

func (l *Listener) onMessage(ctx context.Context, message kafka.Message) error {
	l.MessageChannel <- message
	return nil
}

func (l *Listener) processMessage(ctx context.Context, msg kafka.Message) error {
	// Create context
	ctxWithAudit := zapctx.WithFields(ctx, zap.Bool("audit", true), zap.String("schema_key", l.SchemaKey))
	startTime := time.Now()

	// Decode message
	l.MetricsClient.Count("ingress_topic_counter", 1, []string{l.TopicName})
	messageDecoded, valid := l.SchemaHandlerClient.Decode(ctx, msg.Value, l.SchemaKey)
	if !valid {
		zapctx.Warn(ctxWithAudit, "[Audit] Message sent to quarantine due to failed decode.")
		l.quarantineHandler(ctxWithAudit, msg.Value)
		l.MetricsClient.Count("egress_topic_counter", 1, []string{l.TopicName, "failed"})
		l.MetricsClient.Histogram("transaction_validation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{l.TopicName, "failed"})

		return nil
	}

	// Add relevant message fields to audit log context
	if messageMap, ok := messageDecoded.(map[string]interface{}); ok {
		// Check if field exists in message
		//   requires entering the first level of the map (which is given by the name of the schema)
		for _, messageValue := range messageMap {
			messageValueMap, ok := messageValue.(map[string]interface{})
			if !ok {
				continue
			}
			for _, auditLogField := range l.AuditLogFields {
				if v, ok := messageValueMap[auditLogField]; ok {
					fieldName, fieldValue := auditLogField, fmt.Sprint(v)
					ctxWithAudit = zapctx.WithFields(ctxWithAudit, zap.String(fieldName, fieldValue))
				}
			}
		}
	}

	// Apply business validation
	if l.SkipValidation {
		zapctx.Info(ctx, "[processorHandler] Skipping message validation.")
	} else if validationOk, err := l.BusinessValidator.Validate(messageDecoded); !validationOk {
		zapctx.Warn(ctxWithAudit, "[Audit] Message sent to quarantine due to failed business validation.", zap.Error(err))
		l.quarantineHandler(ctxWithAudit, msg.Value)
		l.MetricsClient.Count("egress_topic_counter", 1, []string{l.TopicName, "failed"})
		l.MetricsClient.Histogram("transaction_validation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{l.TopicName, "failed"})
		return nil
	}

	// Build egress message
	newMessage := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
	}

	// Publish egress message
	if err := l.Producer.WriteMessage(ctx, newMessage); err != nil {
		zapctx.Warn(ctx, "[processorHandler] Failed to write message. Sending to quarantine", zap.Error(err))
		l.quarantineHandler(ctxWithAudit, msg.Value)
		l.MetricsClient.Count("egress_topic_counter", 1, []string{l.TopicName, "failed"})
		l.MetricsClient.Histogram("transaction_validation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{l.TopicName, "failed"})
		return err
	}

	zapctx.Info(ctx, fmt.Sprintf("[processorHandler] Message published to validated topic. Offset: %v", msg.Offset))
	l.MetricsClient.Count("egress_topic_counter", 1, []string{l.TopicName, "valid"})
	l.MetricsClient.Histogram("transaction_validation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{l.TopicName, "valid"})

	return nil
}
