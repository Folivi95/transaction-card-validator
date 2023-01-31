package s3

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	Client *s3.S3
	Bucket string
}

func NewS3Client(awsConfig *aws.Config, bucket string) (*Client, error) {
	session, err := session.NewSession(awsConfig)
	if err != nil {
		return &Client{}, err
	}

	svc := s3.New(session)

	return &Client{
		Client: svc,
		Bucket: bucket,
	}, nil
}

func (client *Client) UploadObject(topicName string, data []byte) error {
	// Write key in hive-like format with Now() timestamp
	year, month, day := time.Now().Date()
	unix := time.Now().UnixMilli()

	key := fmt.Sprintf("%s/year=%d/month=%d/day=%d/%d.json", topicName, year, month, day, unix)

	_, err := client.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(client.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	return err
}
