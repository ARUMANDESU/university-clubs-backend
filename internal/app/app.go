package app

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/app/httpsvr"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/awsS3"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/post"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler"
	clubhandler "github.com/ARUMANDESU/university-clubs-backend/internal/handler/club"
	eventhandler "github.com/ARUMANDESU/university-clubs-backend/internal/handler/event"
	posthandler "github.com/ARUMANDESU/university-clubs-backend/internal/handler/post"
	userhandler "github.com/ARUMANDESU/university-clubs-backend/internal/handler/user"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log/slog"
)

type App struct {
	HTTPSvr *httpsvr.Server
}

// New initializes and returns a new instance of the App struct.
// This function is responsible for the setup and initialization of the core components
// of the application, including the user service client, handlers, and the HTTP server.
// It sets up the necessary dependencies and configurations needed for the application to run.
//
// Parameters:
//   - ctx: A context.Context used to control the lifetime of the user service client.
//   - cfg: A pointer to the config.Config struct containing configuration settings for the application.
//   - log: A *slog.Logger for logging messages and errors throughout the application.
//
// Returns:
//   - A pointer to an initialized App struct, which contains the HTTP server ready to handle requests.
//
// Error Handling:
//   - If the initialization of the user service client fails, the function logs the error and
//     terminates the application using panic. This is typically indicative of a critical error
//     where the application cannot function correctly.
//
// Usage:
//   - This function is usually called at the start of the main function to set up the application.
//     After calling this function, the HTTP server can be started to begin handling requests.
func New(ctx context.Context, cfg *config.Config, log *slog.Logger, awsCfg aws.Config) *App {
	userClient, err := user.New(ctx, log, cfg.Clients.User.Address, cfg.Clients.User.Timeout, cfg.Clients.User.RetriesCount)
	if err != nil {
		log.Error("user service client init error", logger.Err(err))
		panic(err)
	}

	clubClient, err := club.New(ctx, log, cfg.Clients.Club.Address, cfg.Clients.Club.Timeout, cfg.Clients.Club.RetriesCount)
	if err != nil {
		log.Error("club service client init error", logger.Err(err))
		panic(err)
	}

	postsClient, err := post.New(ctx, log, cfg.Clients.Event.Address, cfg.Clients.Event.Timeout, cfg.Clients.Event.RetriesCount)
	if err != nil {
		log.Error("event service client init error", logger.Err(err))
		panic(err)
	}

	cred, err := confidential.NewCredFromSecret(cfg.MicrosoftOIDC.Secret)
	if err != nil {
		log.Error("failed to create a Credential from a secret.", logger.Err(err))
		panic(err)
	}
	confidentialClient, err := confidential.New(cfg.MicrosoftOIDC.Authority, cfg.MicrosoftOIDC.ClientID, cred)
	if err != nil {
		log.Error("failed to create a oidc client.", logger.Err(err))
		panic(err)
	}

	awsS3Storage, err := awsS3.New(awsCfg)
	if err != nil {
		log.Error("failed to create aws s3 client", logger.Err(err))
		panic(err)
	}

	h := handler.New(cfg, handler.Handlers{
		UsrHandler:   userhandler.New(cfg, log, userClient, confidentialClient, awsS3Storage),
		ClubHandler:  clubhandler.New(log, clubClient, awsS3Storage),
		EventHandler: eventhandler.New(log, postsClient, clubClient, userClient, awsS3Storage, awsS3Storage),
		PostHandler:  posthandler.New(log, postsClient),
	})

	httpServer := httpsvr.New(cfg, h.InitRoutes())

	return &App{HTTPSvr: httpServer}
}
