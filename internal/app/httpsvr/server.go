package httpsvr

import (
	"fmt"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"net/http"
)

type Server struct {
	HTTPServer *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	httpServer := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &Server{HTTPServer: httpServer}
}

func (s Server) MustRun() {
	err := s.Run()
	if err != nil {
		panic(err)
	}
}

func (s Server) Run() error {
	const op = "app.Run"
	err := s.HTTPServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Server) Stop() {

}
