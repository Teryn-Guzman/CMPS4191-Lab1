package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startServer(app *applicationDependencies, handler http.Handler) error {
    addr := fmt.Sprintf(":%d", app.config.port)
    srv := &http.Server{
        Addr:    addr,
        Handler: handler,
    }

    go func() {
        app.logger.Info("starting server", "addr", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            app.logger.Error("listen failed", "err", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    app.logger.Info("shutdown signal received, shutting down server")
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.config.shutdownTimeout)*time.Second)
    defer cancel()
    return srv.Shutdown(ctx)
}
