package club

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	"log/slog"
)

type Handler struct {
	clubClient   *club.Client
	log          *slog.Logger
	imageStorage ImageStorage
}

type ImageStorage interface {
	UploadImage(ctx context.Context, image []byte, filename string, bucket string) (string, error)
	DeleteImage(ctx context.Context, filename string, bucket string) error
}

// New creates and returns a new User Handler instance
// Parameters:
//   - client: A *user.Client which is a gRPC client for the user service.
//   - log: A *slog.Logger used for logging messages and errors.
//
// Returns:
//   - A Handler struct that encapsulates the provided user service client and logger.
func New(log *slog.Logger, client *club.Client, imageStorage ImageStorage) Handler {
	return Handler{
		clubClient:   client,
		log:          log,
		imageStorage: imageStorage,
	}
}
