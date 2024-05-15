package events

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/awsS3"
)

type Handler struct {
	s3Client     *awsS3.Client
	imageStorage ImageStorage
}

type ImageStorage interface {
	UploadImage(ctx context.Context, image []byte, filename string, bucket string) (string, error)
	DeleteImage(ctx context.Context, filename string, bucket string) error
}

func New(client *awsS3.Client, imageStorage ImageStorage) Handler {
	return Handler{
		s3Client:     client,
		imageStorage: imageStorage,
	}
}
