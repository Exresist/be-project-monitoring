package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"be-project-monitoring/internal/domain/model"
)

type (
	server struct {
		*http.Server
		logger   *zap.Logger
		response *responder
		svc      service

		shutdownTimeout int
	}

	service interface {
		VerifyToken(ctx context.Context, token string, toAllow ...model.UserRole) error
		CreateUser(ctx context.Context, user *model.User) (*model.User, string, error)
	}

	OptionFunc func(s *server)
)

func New(opts ...OptionFunc) *server {
	s := &server{}
	for _, opt := range opts {
		opt(s)
	}
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TokenClaims"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/register", s.createUser)

	r.Route("/admin", func(r chi.Router) {
		r.Use(s.authMiddleware(model.Admin))
		r.Get("/users", func(writer http.ResponseWriter, request *http.Request) {
		})
		r.Get("/projects", func(writer http.ResponseWriter, request *http.Request) {
		})
	})

	s.Handler = r
	return s
}

func (s *server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err := s.ListenAndServe()

		if err != nil && errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("http-server failed to launch: %w", err)
	})

	eg.Go(func() error {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			time.Duration(s.shutdownTimeout)*time.Second,
		)
		defer shutdownCancel()

		//nolint:contextcheck
		if err := s.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http-server shutted down: %w", err)
		}

		return nil
	})
	return eg.Wait()
}

func WithLogger(logger *zap.Logger) OptionFunc {
	return func(s *server) {
		s.logger = logger
	}
}

func WithServer(srv *http.Server) OptionFunc {
	return func(s *server) {
		s.Server = srv
	}
}

func WithResponder(r *responder) OptionFunc {
	return func(s *server) {
		s.response = r
	}
}

func WithService(svc service) OptionFunc {
	return func(s *server) {
		s.svc = svc
	}
}

func WithShutdownTimeout(timeout int) OptionFunc {
	return func(s *server) {
		s.shutdownTimeout = timeout
	}
}
