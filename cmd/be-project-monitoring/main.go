package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
	"go.uber.org/zap"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/repository"
	"be-project-monitoring/internal/domain/service"
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

	/*logger.Info("creating HTTP server")
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Duration(cfg.WriteTimeout),
	}*/
	var g = &run.Group{}

	userStore := repository.NewUserStore(conn, "users", sugaredLogger)
	userSvc := service.NewUserService(userStore)
	projectStore := repository.NewProjectStore(conn, "projects", sugaredLogger)
	projSvc := service.NewProjectService(projectStore)
	
	api.New(
		// api.WithServer(srv),
		api.WithLogger(sugaredLogger),
		api.WithUserService(userSvc),
		api.WithProjectService(projSvc),
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
	//change
}
