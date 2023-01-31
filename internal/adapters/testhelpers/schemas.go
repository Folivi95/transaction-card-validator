package testhelpers

import (
	"context"
	"errors"

	"github.com/linkedin/goavro/v2"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"

	solar "github.com/saltpay/transaction-card-validator/internal/adapters/testhelpers/solanteq"
	w4 "github.com/saltpay/transaction-card-validator/internal/adapters/testhelpers/way4"
)

type LocalModel struct{}

func (l *LocalModel) getSchema(schemaKey string) (string, error) {
	switch schemaKey {
	case "Way4RawTransaction":
		return w4.DocRawTransaction, nil
	case "Way4CuratedTransaction":
		return w4.DocTokenizedTransaction, nil
	case "SolanteqRawTransaction":
		return solar.AfterProcessTxnEnvelope, nil
	case "SolanteqRawPayoutInstruction":
		return solar.AfterInvoiceIssuingEnvelope, nil
	default:
		return "", errors.New("schema definition not founded in models")
	}
}

func (l *LocalModel) GetNewCodec(ctx context.Context, schemaKey string) (goavro.Codec, error) {
	schemaStr, err := l.getSchema(schemaKey)
	if err != nil {
		return goavro.Codec{}, err
	}

	codec, err := goavro.NewCodecForStandardJSON(schemaStr)
	if err != nil {
		return goavro.Codec{}, err
	}

	return *codec, nil
}

func (l *LocalModel) Decode(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
	// Schema validation
	codec, err := l.GetNewCodec(ctx, schemaKey)
	if err != nil {
		zapctx.Error(ctx, "[processorHandler] Failed to create Codec for given schema. Sending message to quarantine", zap.Error(err))
		return nil, false
	}

	nativeType, _, err := codec.NativeFromTextual(msg)
	if err != nil {
		zapctx.Warn(ctx, "[Audit] Message does not comply to the defined schema. Sending it to quarantine with omitted CHD.")
		return nil, false
	}

	return nativeType, true
}

func NewLocalModel() *LocalModel {
	return &LocalModel{}
}
