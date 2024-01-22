package app

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/internal/app/httpsvr"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler"
	"log/slog"
)

type App struct {
	HTTPSvr *httpsvr.Server
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger) *App {
	userClient, err := user.New(ctx, log, cfg.Clients.User.Address, cfg.Clients.User.Timeout, cfg.Clients.User.RetriesCount)
	if err != nil {
		log.Error("user service client init error", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		panic(err)
	}

	h := handler.New(userClient, log)

	httpServer := httpsvr.New(cfg, h.InitRoutes())

	return &App{HTTPSvr: httpServer}
}
