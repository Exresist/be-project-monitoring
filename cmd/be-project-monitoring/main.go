package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/repository"
	"be-project-monitoring/internal/domain/service"
	ierr "be-project-monitoring/internal/errors"
)

func main() {
	var (
		ctx, ctxCancel = context.WithCancel(context.Background())
		cfg            = new(config)
		logger         *zap.Logger
		err            error
	)
	defer ctxCancel()

	if err = envconfig.Process("APP", cfg); err != nil {
		log.Fatal(err.Error())
	}

	if cfg.Env == "development" {
		logger, err = zap.NewDevelopment()

	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}

	conn, err := db.ConnectPostgreSQL(ctx, cfg.DSN)
	if err != nil {
		panic(fmt.Errorf("невозможно открыть соединение с базой данных: %w", err))
	}

	logger.Info("creating HTTP server")
	srv := &http.Server{
		Addr:         cfg.BindAddr,
		ReadTimeout:  time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Duration(cfg.WriteTimeout),
	}

	eg, egCtx := errgroup.WithContext(ctx)
	{
		logger.Info("start of the server")
		eg.Go(func() error {
			srvCtx, srvCancel := context.WithCancel(egCtx)
			defer srvCancel()

			userStore := repository.NewUserStore(conn, "users", logger)

			svc := service.NewService(userStore)
			return api.New(
				api.WithServer(srv),
				api.WithLogger(logger),
				api.WithResponder(api.NewResponder(logger)),
				api.WithService(svc),
				api.WithShutdownTimeout(cfg.ShutdownTimeout)).Run(srvCtx)
		})
	}

	eg.Go(func() error {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

		logger.Info("[signal-watcher] started")

		select {
		case sig := <-shutdown:
			return fmt.Errorf("%w: %s", ierr.ErrTermSig, sig.String())
		case <-egCtx.Done():
			return nil
		}
	})

	defer func() {
		recovered := recover()
		if e, ok := recovered.(error); ok && errors.Is(e, ierr.ErrAbnormalExit) {
			os.Exit(1)
		}
	}()

	if err := eg.Wait(); err != nil &&
		!errors.Is(err, ierr.ErrTermSig) &&
		!errors.Is(err, context.Canceled) {
		logger.Error("emergency service shutdown", zap.Error(err))
		panic(ierr.ErrAbnormalExit)
	}
	logger.Info("successful shutdown")
}
