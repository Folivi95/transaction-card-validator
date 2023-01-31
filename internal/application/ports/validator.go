//go:generate moq -out mocks/validator_moq.go -pkg=mocks . Validator

package ports

type Validator interface {
	Validate(msgDecoded interface{}) (bool, error)
}

type DummyValidator struct{}

func (v *DummyValidator) Validate(msgDecoded interface{}) (bool, error) {
	return true, nil
}
