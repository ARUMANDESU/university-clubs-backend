package app

import (
	"github.com/ARUMANDESU/university-clubs-backend/internal/app/httpsvr"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler"
	"log/slog"
)

type App struct {
	log     *slog.Logger
	HTTPSvr *httpsvr.Server
}

func New(cfg *config.Config, log *slog.Logger) *App {
	h := handler.New()

	httpServer := httpsvr.New(cfg, h.InitRoutes())

	return &App{log: log, HTTPSvr: httpServer}
}
