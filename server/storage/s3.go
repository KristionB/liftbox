package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadObject(client *s3.Client, bucket, key string, data []byte) error {
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}

