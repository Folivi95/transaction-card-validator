//go:build unit

package way4validation_test

import (
	"testing"

	"github.com/matryer/is"

	way4validation "github.com/saltpay/transaction-card-validator/internal/application/way4"
)

const validTestMessage = `{"f1": "field1"}`

func TestWay4Validation(t *testing.T) {
	is := is.New(t)
	t.Run("validate a test schema", func(t *testing.T) {
		message := []byte(validTestMessage)

		validator, err := way4validation.NewWay4MaskedCardNumber([]string{}, "bar")
		is.NoErr(err)
		validationOk, err := validator.Validate(message)
		is.True(validationOk)
		is.Equal(err, nil)
	})
}
