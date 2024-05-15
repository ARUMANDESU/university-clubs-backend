package event

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/event"
)

type Handler struct {
	eventClient  *event.Client
	imageStorage ImageStorage
}

type ImageStorage interface {
	UploadImage(ctx context.Context, image []byte, filename string, bucket string) (string, error)
	DeleteImage(ctx context.Context, filename string, bucket string) error
}

func New(client *event.Client, imageStorage ImageStorage) Handler {
	return Handler{
		eventClient:  client,
		imageStorage: imageStorage,
	}
}
