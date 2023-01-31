//go:generate moq -out mocks/refresh_scheduler_moq.go -pkg mocks . RefreshScheduler
//go:generate moq -out mocks/schema_registry_moq.go -pkg mocks . SchemaRegistry

package schemaregistry

import (
	"context"
	"time"

	"github.com/saltpay/transaction-card-validator/internal/adapters/schemaregistry/codecclient"

	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"
)

type RefreshScheduler interface {
	AfterFunc(t time.Duration, f func()) *time.Timer
}

type SchedulerFunc func(d time.Duration, f func()) *time.Timer

func (s SchedulerFunc) AfterFunc(d time.Duration, f func()) *time.Timer {
	return s(d, f)
}

type SchemaRegistry interface {
	GetLatestSchema(subject string) (*srclient.Schema, error)
	GetSchemaByVersion(subject string, version int) (*srclient.Schema, error)
}

type Client struct {
	schemaTable map[string][]string
	Session     SchemaRegistry
	CodecClient codecclient.ICodecClient
	Scheduler   RefreshScheduler
}

func NewSchemaRegistryClient(endpoint string, refreshIntervalSeconds int) *Client {
	schemaTable := make(map[string][]string)
	session := srclient.CreateSchemaRegistryClient(endpoint)
	scheduler := SchedulerFunc(time.AfterFunc)

	client := &Client{
		schemaTable: schemaTable,
		Session:     session,
		Scheduler:   scheduler,
		CodecClient: &codecclient.CodecClient{
			CodecTable: map[string]*goavro.Codec{},
		},
	}
	client.RefreshSchemasTrigger(refreshIntervalSeconds)

	return client
}

func (c *Client) Decode(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
	schema, err := c.getLatestSchema(ctx, schemaKey)
	if err != nil {
		zapctx.Error(ctx, "[Client] Unable to fetch schema to create codec", zap.Error(err))
	}
	return c.CodecClient.Decode(ctx, msg, schemaKey, schema)
}

func (c *Client) GetNewCodec(ctx context.Context, schemaKey string) (*goavro.Codec, error) {
	schema, err := c.getLatestSchema(ctx, schemaKey)
	if err != nil {
		zapctx.Error(ctx, "[Client] Unable to fetch schema to create codec", zap.Error(err))
	}
	return c.CodecClient.GetNewCodec(ctx, schemaKey, schema)
}

func (c *Client) getLatestSchema(ctx context.Context, schemaKey string) ([]string, error) {
	schema, schemaRegistered := c.schemaTable[schemaKey]

	var err error
	if !schemaRegistered {
		schema, err = c.fetchSchema(ctx, schemaKey)
		if err != nil {
			return []string{}, err
		}
	}

	return schema, nil
}

func (c *Client) fetchSchema(ctx context.Context, schemaKey string) ([]string, error) {
	schema, err := c.Session.GetLatestSchema(schemaKey)
	if err != nil {
		return []string{}, err
	}

	var newSchema []string

	// Handle references
	if schema.References() != nil {
		// Build list of all schema references
		//  by iterating all references (and children references) in a BFS manner
		schemaReferences := make([]srclient.Schema, 0, 10)
		schemaReferencesIterationStack := schema.References()
		for len(schemaReferencesIterationStack) > 0 {
			// Pop first element
			reference := schemaReferencesIterationStack[0]
			schemaReferencesIterationStack = schemaReferencesIterationStack[1:]

			// Get schema from registry
			refSchema, err := c.Session.GetSchemaByVersion(reference.Subject, reference.Version)
			if err != nil {
				zapctx.Error(ctx, "[registryHandler] Failed to get reference schema", zap.Error(err))
				continue
			}

			// Append result
			schemaReferences = append([]srclient.Schema{*refSchema}, schemaReferences...)

			// Iterate new references
			schemas := refSchema.References()
			schemaReferencesIterationStack = append(schemaReferencesIterationStack, schemas...)
		}

		// Remove duplicates from schemaReferences (duplicates may happen if multiple schemas have a common reference)
		schemasVisited := make(map[int]bool)
		schemasToDeleteIdx := make([]int, 0, 10)
		for i, schema := range schemaReferences {
			if _, exists := schemasVisited[schema.ID()]; exists {
				// Mark schema to be removed from schema references (duplicate)
				schemasToDeleteIdx = append(schemasToDeleteIdx, i)
			} else {
				// Keep track that this schema exists
				schemasVisited[schema.ID()] = true
			}
		}
		// Iterate schemasToDeleteIdx in reverser order to delete elements
		for i := len(schemasToDeleteIdx) - 1; i >= 0; i-- {
			schemaIdxToRemove := schemasToDeleteIdx[i]
			schemaReferences = append(schemaReferences[:schemaIdxToRemove], schemaReferences[schemaIdxToRemove+1:]...)
		}

		// Dump all references to newSchema variable
		for _, reference := range schemaReferences {
			newSchema = append(newSchema, reference.Schema())
		}
	}

	// Appends envelope schema to entities schema
	newSchema = append(newSchema, schema.Schema())

	if len(newSchema) == 0 {
		zapctx.Error(ctx, "[registryHandler] Empty schema for given subject name. Aborting.", zap.Error(err))
		return []string{}, err
	}

	c.schemaTable[schemaKey] = newSchema
	return newSchema, err
}

func (c *Client) RefreshSchemasTrigger(refreshIntervalSeconds int) {
	c.refreshSchemas(refreshIntervalSeconds)
}

func (c *Client) refreshSchemas(refreshIntervalSeconds int) {
	ctx := context.Background()
	for schemaKey := range c.schemaTable {
		newSchema, err := c.fetchSchema(ctx, schemaKey)
		if err != nil {
			zapctx.Error(ctx, "[refreshSchemas] Failed to retrieve schema", zap.Error(err))
			continue
		}
		c.schemaTable[schemaKey] = newSchema
	}

	c.Scheduler.AfterFunc(time.Duration(refreshIntervalSeconds)*time.Second, func() {
		c.RefreshSchemasTrigger(refreshIntervalSeconds)
		c.CodecClient.RefreshCodec()
	})
}
