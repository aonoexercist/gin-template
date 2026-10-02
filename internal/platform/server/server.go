package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/yourname/gin-template/internal/platform/config"
)

type Server struct {
	http    *http.Server
	log     *slog.Logger
	timeout time.Duration
}

func New(cfg *config.Config, log *slog.Logger, h http.Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           h,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		log:     log,
		timeout: cfg.ShutdownTimeout,
	}
}

// Run blocks until ctx is cancelled, then shuts down gracefully.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("server starting", "addr", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	s.log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.http.Shutdown(shutdownCtx)
}
