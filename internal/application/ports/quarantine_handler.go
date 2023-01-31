//go:generate moq -out mocks/quarantine_handler_moq.go -pkg=mocks . QuarantineHandler

package ports

type QuarantineHandler interface {
	UploadObject(topicName string, data []byte) error
}
