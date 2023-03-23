package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/service"
	"be-project-monitoring/internal/repository"

	"github.com/google/go-github/v49/github"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func main() {
	var (
		ctx, ctxCancel = context.WithCancel(context.Background())
		cfg            = new(config)
		logger         *zap.Logger
		err            error
	)
	defer ctxCancel()

	if cfg.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	sugaredLogger := logger.Sugar()

	if err = envconfig.Process("APP", cfg); err != nil {
		sugaredLogger.Fatal(err.Error())
	}

	conn, err := db.ConnectPostgreSQL(ctx, cfg.DSN)
	if err != nil {
		panic(fmt.Errorf("невозможно открыть соединение с базой данных: %w", err))
	}

	var g = &run.Group{}

	repo := repository.NewRepository(conn, sugaredLogger)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GHTOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubCl := github.NewClient(tc)
	svc := service.NewService(repo, githubCl)
	api.New(
		api.WithLogger(sugaredLogger),
		api.WithService(svc),
		api.WithShutdownTimeout(cfg.ShutdownTimeout)).Run(g)

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(func() error {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

		sugaredLogger.Info("[signal-watcher] started")

		select {
		case sig := <-shutdown:
			return fmt.Errorf("terminated with signal: %s", sig.String())
		case <-ctx.Done():
			return nil
		}
	}, func(err error) {
		cancel()
		sugaredLogger.Error("gracefully shutdown application", zap.Error(err))
	})

	sugaredLogger.Error("successful shutdown", zap.Error(g.Run()))
}
