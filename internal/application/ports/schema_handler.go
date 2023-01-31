//go:generate moq -out mocks/schema_handler_moq.go -pkg=mocks . SchemaHandler

package ports

import (
	"context"
)

type SchemaHandler interface {
	Decode(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool)
}
