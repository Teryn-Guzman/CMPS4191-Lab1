package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (a *applicationDependencies) serve() error {
	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.port),
		Handler:      routes(a),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		s := <-quit
		a.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(a.config.shutdownTimeout)*time.Second,
		)
		defer cancel()

		shutdownErr := apiServer.Shutdown(ctx)
		a.logger.Info("completing background tasks", "address", apiServer.Addr)
		a.wg.Wait()
		shutdownError <- shutdownErr
	}()

	a.logger.Info("starting server", "address", apiServer.Addr, "environment", a.config.environment)
	err := apiServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err = <-shutdownError; err != nil {
		return err
	}

	a.logger.Info("stopped server", "address", apiServer.Addr)
	return nil
}
