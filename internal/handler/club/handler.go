package club

import (
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	"log/slog"
)

type Handler struct {
	clubClient *club.Client
	log        *slog.Logger
}

// New creates and returns a new User Handler instance
// Parameters:
//   - client: A *user.Client which is a gRPC client for the user service.
//   - log: A *slog.Logger used for logging messages and errors.
//
// Returns:
//   - A Handler struct that encapsulates the provided user service client and logger.
func New(client *club.Client, log *slog.Logger) Handler {
	return Handler{
		clubClient: client,
		log:        log,
	}
}
