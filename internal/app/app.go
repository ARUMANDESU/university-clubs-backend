package app

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/app/httpsvr"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/notification"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler"
	"github.com/ARUMANDESU/university-clubs-backend/internal/services/notification"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	socketio "github.com/googollee/go-socket.io"
	"log/slog"
)

type App struct {
	HTTPSvr             *httpsvr.Server
	NotificationService *notificationSrvc.Service
	SocketIOSrv         *socketio.Server
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
func New(ctx context.Context, cfg *config.Config, log *slog.Logger) *App {
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

	//todo: take from config
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

	notificationClient, err := notification.New(log, cfg.Clients.Notification.Address, cfg.Clients.Notification.Timeout, cfg.Clients.Notification.RetriesCount)
	if err != nil {
		log.Error("notification service client init error", logger.Err(err))
		panic(err)
	}

	notificationService := notificationSrvc.NewService(log, notificationClient)

	h := handler.New(log, cfg.MicrosoftOIDC, userClient, clubClient, confidentialClient, notificationService)

	routes, socketIOSrv := h.InitRoutes()

	httpServer := httpsvr.New(cfg, routes)

	return &App{HTTPSvr: httpServer, NotificationService: notificationService, SocketIOSrv: socketIOSrv}
}
