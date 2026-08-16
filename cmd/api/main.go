package main

import (
	"context"
	"expvar"
	"flag"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/Teryn-Guzman/Gatekeeper/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

const appVersion = "0.1.0"

type serverConfig struct {
    port int
    environment string
    db struct{
        dsn string
    }
    shutdownTimeout int
}

type applicationDependencies struct{
    config serverConfig
    logger *slog.Logger
    store  *store.Store
    wg     sync.WaitGroup
}

func main() {
    var cfg serverConfig

    flag.IntVar(&cfg.port, "port", 8080, "Server port")
    flag.StringVar(&cfg.environment, "env", "development", "Environment (development|staging|production)")
    flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "Postgres DSN")
    flag.IntVar(&cfg.shutdownTimeout, "shutdown-timeout", 15, "Graceful shutdown timeout in seconds")
    flag.Parse()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    if cfg.db.dsn == "" {
        logger.Error("database DSN is required via -db-dsn or DATABASE_URL env")
        os.Exit(1)
    }

    ctx := context.Background()
    pool, err := pgxpool.New(ctx, cfg.db.dsn)
    if err != nil {
        logger.Error("failed to connect to db", "err", err)
        os.Exit(1)
    }
    defer pool.Close()

    expvar.NewString("version").Set(appVersion)
    expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
    expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))

    app := &applicationDependencies{
        config: cfg,
        logger: logger,
        store:  store.New(pool),
    }

    if err := app.serve(); err != nil {
        logger.Error("server error", "err", err)
        os.Exit(1)
    }
}

func (app *applicationDependencies) serve() error {
    handler := routes(app)
    return startServer(app, handler)
}

