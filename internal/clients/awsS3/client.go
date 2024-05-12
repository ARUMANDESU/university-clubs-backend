package awsS3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client   *s3.Client
	uploader *manager.Uploader
}

func New(awsCfg aws.Config) (*Client, error) {
	client := s3.NewFromConfig(awsCfg)
	uploader := manager.NewUploader(client)

	return &Client{
		client:   client,
		uploader: uploader,
	}, nil
}

func (c *Client) UploadImage(ctx context.Context, image []byte, filename string, bucket string) (string, error) {

	result, err := c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(image),
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (c *Client) DeleteImage(ctx context.Context, filename string, bucket string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}

	return nil
}
