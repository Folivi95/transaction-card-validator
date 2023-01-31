//go:build unit

package codecclient_test

import (
	"context"
	"github.com/matryer/is"
	"github.com/riferrei/srclient"
	"testing"
	"time"

	"github.com/saltpay/transaction-card-validator/internal/adapters/schemaregistry"
	"github.com/saltpay/transaction-card-validator/internal/adapters/schemaregistry/mocks"
)

const validSchema = `{
	"type": "record",
	"name": "name",
	"doc": "Some data",
	"fields": [
		{"name": "f1", "type": "string"}
	]
}`

var testSchema = make([]string, 3)

func TestCodecClient_Decode(t *testing.T) {
	is := is.New(t)

	t.Run("client should fetch codec and decode valid messages", func(t *testing.T) {
		testSchema = append(testSchema, validSchema)
		validMessage := `{"f1":"field1"}`
		expectedResponse := "field1"

		registryClient := schemaregistry.NewSchemaRegistryClient("some_endpoint", 5)
		mockedRegistry := &mocks.SchemaRegistryMock{
			GetLatestSchemaFunc: func(subject string) (*srclient.Schema, error) {
				schema, _ := srclient.NewSchema(1, validSchema, "PROTOBUF", 2, nil, nil, nil)
				return schema, nil
			},
		}
		registryClient.Session = mockedRegistry
		_, _ = registryClient.Decode(context.TODO(), []byte(validMessage), "some_key")
		messageInterface, valid := registryClient.CodecClient.Decode(context.TODO(),
			[]byte(validMessage),
			"some_key",
			testSchema)

		is.True(valid)
		is.Equal(messageInterface.(map[string]interface{})["name"].(map[string]interface{})["f1"], expectedResponse)
		is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 1)
	})
}

func TestCodecClient_RefreshCodec(t *testing.T) {
	is := is.New(t)

	t.Run("client should refresh codec cache after specified time", func(t *testing.T) {
		testSchema = append(testSchema, validSchema)
		authCalled := make(chan bool)
		registryClient := schemaregistry.NewSchemaRegistryClient("some_endpoint", 1)
		mockedRegistry := &mocks.SchemaRegistryMock{
			GetLatestSchemaFunc: func(subject string) (*srclient.Schema, error) {
				schema, _ := srclient.NewSchema(1, validSchema, "PROTOBUF", 2, nil, nil, nil)
				return schema, nil
			},
		}
		mockedScheduler := &mocks.RefreshSchedulerMock{
			AfterFuncFunc: func(t time.Duration, f func()) *time.Timer {
				go func() {
					authCalled <- true
				}()
				return time.NewTimer(2 * time.Second)
			},
		}
		registryClient.Session = mockedRegistry
		registryClient.Scheduler = mockedScheduler

		_, err := registryClient.GetNewCodec(context.TODO(), "some_key")
		is.NoErr(err)
		is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 1)

		select {
		case <-authCalled:
			is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 2)
		case <-time.After(2000 * time.Millisecond):
			t.Fatal("timed out before refreshing tokens")
		}

	})

	t.Run("size of codec cache after refresh should be zero", func(t *testing.T) {
		testSchema = append(testSchema, validSchema)
		authCalled := make(chan bool)
		registryClient := schemaregistry.NewSchemaRegistryClient("some_endpoint", 1)
		mockedRegistry := &mocks.SchemaRegistryMock{
			GetLatestSchemaFunc: func(subject string) (*srclient.Schema, error) {
				schema, _ := srclient.NewSchema(1, validSchema, "PROTOBUF", 2, nil, nil, nil)
				return schema, nil
			},
		}
		mockedScheduler := &mocks.RefreshSchedulerMock{
			AfterFuncFunc: func(t time.Duration, f func()) *time.Timer {
				go func() {
					authCalled <- true
				}()
				return time.NewTimer(2 * time.Second)
			},
		}
		registryClient.Session = mockedRegistry
		registryClient.Scheduler = mockedScheduler

		_, err := registryClient.GetNewCodec(context.TODO(), "some_key")
		is.NoErr(err)
		is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 1)

		select {
		case <-authCalled:
			is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 2)
		case <-time.After(2000 * time.Millisecond):
			t.Fatal("timed out before refreshing tokens")
		}

		is.Equal(registryClient.CodecClient.Size(), 0)
	})
}
