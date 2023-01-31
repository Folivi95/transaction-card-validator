package codecclient

import (
	"context"
	"strings"
	"sync"

	"github.com/linkedin/goavro/v2"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"
)

type ICodecClient interface {
	Decode(ctx context.Context, msg []byte, schemaKey string, schema []string) (interface{}, bool)
	GetNewCodec(ctx context.Context, schemaKey string, schema []string) (*goavro.Codec, error)
	RefreshCodec()
	Size() int
}

type CodecClient struct {
	sync.Mutex
	CodecTable map[string]*goavro.Codec
}

func (c *CodecClient) Decode(ctx context.Context, msg []byte, schemaKey string, schema []string) (interface{}, bool) {
	codec, err := c.GetNewCodec(ctx, schemaKey, schema)
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

func (c *CodecClient) GetNewCodec(ctx context.Context, schemaKey string, schema []string) (*goavro.Codec, error) {
	var err error
	c.Lock()
	defer c.Unlock()
	codec, codecRegistered := c.CodecTable[schemaKey]

	if !codecRegistered {
		// Convert schema slice to string in Avro format
		schemaStr := "[" + strings.Join(schema, ",") + "]"

		codec, err := goavro.NewCodecForStandardJSON(schemaStr)
		if err != nil {
			zapctx.Error(ctx, "[registryHandler] Failed to create Codec for given schema.", zap.Error(err))
			return nil, err
		}
		c.CodecTable[schemaKey] = codec
		return codec, err
	}
	return codec, err
}

func (c *CodecClient) RefreshCodec() {
	c.CodecTable = make(map[string]*goavro.Codec)
}

func (c *CodecClient) Size() int {
	return len(c.CodecTable)
}
