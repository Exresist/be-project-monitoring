package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"be-project-monitoring/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"go.uber.org/zap"
)

type (
	Server struct {
		*http.Server
		logger *zap.SugaredLogger
		svc    Service

		shutdownTimeout int
	}

	Service interface {
		userService
		projectService
		participantService
	}
	userService interface {
		VerifyToken(ctx context.Context, token string, toAllow ...model.UserRole) error
		CreateUser(ctx context.Context, user *CreateUserReq) (*model.User, string, error)
		AuthUser(ctx context.Context, username, password string) (string, error)
		GetUsers(ctx context.Context, userReq *GetUserReq) ([]model.User, int, error)
		FindGithubUser(ctx context.Context, userReq string) bool
		UpdateUser(ctx context.Context, userReq *UpdateUserReq) (*model.User, error)
		DeleteUser(ctx context.Context, userReq *DeleteUserReq) error
	}

	projectService interface {
		CreateProject(ctx context.Context, project *CreateProjectReq) (*model.Project, error)
		UpdateProject(ctx context.Context, projectReq *UpdateProjectReq) (*model.Project, error)
		DeleteProject(ctx context.Context, projectReq *DeleteProjectReq) error
		GetProjects(ctx context.Context, projectReq *GetProjectReq) ([]model.Project, int, error)
	}

	participantService interface {
		AddParticipant(ctx context.Context, participant *model.Participant) ([]model.Participant, error)
		GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error)
	}

	OptionFunc func(s *Server)
)

func New(opts ...OptionFunc) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:         ":8080",
			ReadTimeout:  time.Duration(10) * time.Second,
			WriteTimeout: time.Duration(10) * time.Second},
	}
	for _, opt := range opts {
		opt(s)
	}

	rtr := gin.Default()

	// /api/*
	apiRtr := rtr.Group("/api")
	// /api/auth
	apiRtr.POST("/auth", s.auth)
	// /api/register
	apiRtr.POST("/register", s.register)

	// /api/pm
	pmRtr := apiRtr.Group("/pm", s.authMiddleware(model.ProjectManager))

	// api/pm/project
	projectRtr := pmRtr.Group("/project")

	projectRtr.PUT("/", s.createProject)
	projectRtr.POST("/:id", s.addParticipant)

	// /api/admin
	adminRtr := apiRtr.Group("/admin", s.authMiddleware(model.Admin))

	// /api/admin/users
	// TODO
	adminRtr.GET("/users", s.getUsers)
	// /api/admin/projects
	adminRtr.GET("/projects", s.getProjects)

	s.Handler = rtr
	return s
}

func (s *Server) Run(g *run.Group) {
	g.Add(func() error {
		s.logger.Info("[http-server] started")
		s.logger.Info(fmt.Sprintf("listening on %v", s.Addr))
		return s.ListenAndServe()
	}, func(err error) {
		s.logger.Error("[http-server] terminated", zap.Error(err))

		ctx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()

		s.logger.Error("[http-server] stopped", zap.Error(s.Shutdown(ctx)))
	})
}

func WithLogger(logger *zap.SugaredLogger) OptionFunc {
	return func(s *Server) {
		s.logger = logger
	}
}

/*func WithServer(srv *http.Server) OptionFunc {
	return func(s *Server) {
		s.Server = srv
	}
}*/

func WithService(svc Service) OptionFunc {
	return func(s *Server) {
		s.svc = svc
	}
}

func WithShutdownTimeout(timeout int) OptionFunc {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
