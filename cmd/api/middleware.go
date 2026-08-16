package main

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            logger.Info("request start", "method", r.Method, "path", r.URL.Path)
            next.ServeHTTP(w, r)
            logger.Info("request end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
        })
    }
}
