package main

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/university-clubs-backend/internal/app"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTPServer.Address),
	)

	ctx := context.Background()

	application := app.New(ctx, cfg, log)

	go func() {
		if err := application.HTTPSvr.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error: %v", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		}
		log.Info("stopped serving new connections")

	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("shutting down application", slog.String("signal", sign.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := application.HTTPSvr.Stop(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("shutdown error", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
	log.Info("graceful shutdown complete")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
