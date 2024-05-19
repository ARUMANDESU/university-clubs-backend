package event

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/event"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"log/slog"
)

type Handler struct {
	eventClient  *event.Client
	clubClient   *club.Client
	userClient   *user.Client
	log          *slog.Logger
	imageStorage ImageStorage
	FileStorage  FileStorage
}

type ImageStorage interface {
	Upload(ctx context.Context, image []byte, filename string, bucket string) (string, error)
	Delete(ctx context.Context, filename string, bucket string) error
}

type FileStorage interface {
	Upload(ctx context.Context, file []byte, filename string, bucket string) (string, error)
	Delete(ctx context.Context, filename string, bucket string) error
}

func New(
	log *slog.Logger,
	client *event.Client,
	clubClient *club.Client,
	userClient *user.Client,
	imageStorage ImageStorage,
	fileStorage FileStorage,
) Handler {
	return Handler{
		eventClient:  client,
		clubClient:   clubClient,
		userClient:   userClient,
		log:          log,
		imageStorage: imageStorage,
		FileStorage:  fileStorage,
	}
}
