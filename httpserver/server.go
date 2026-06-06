package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Options struct {
	Server          *http.Server
	Logger          *slog.Logger
	ShutdownTimeout time.Duration
	Close           func() error
}

func Run(ctx context.Context, opts Options) error {
	if opts.Server == nil {
		return fmt.Errorf("http server is required")
	}
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	shutdownTimeout := opts.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 10 * time.Second
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("http server starting", slog.String("http_addr", opts.Server.Addr))
		if err := opts.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := opts.Server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}
		if opts.Close != nil {
			if err := opts.Close(); err != nil {
				return fmt.Errorf("close app: %w", err)
			}
		}
		log.Info("http server stopped")
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("run http server: %w", err)
		}
		return nil
	}
}
