package api

import (
	"be-project-monitoring/internal/domain/model"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		taskService
		tokenService
	}
	userService interface {
		CreateUser(ctx context.Context, user *CreateUserReq) (*model.User, string, error)
		AuthUser(ctx context.Context, username, password string) (*model.User, string, error)
		GetFullUsers(ctx context.Context, searchParam string) ([]model.User, int, error)
		GetPartialUsers(ctx context.Context, userReq *GetUserReq) ([]model.ShortUser, int, error)
		FindGithubUser(ctx context.Context, userReq string) bool
		UpdateUser(ctx context.Context, userReq *UpdateUserReq) (*model.User, error)
		DeleteUser(ctx context.Context, id uuid.UUID) error
		GetUserProfile(ctx context.Context, id uuid.UUID) (*model.Profile, error)
	}

	tokenService interface {
		GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error)
		VerifyToken(ctx context.Context, token string, toAllow ...model.UserRole) error
		VerifySelf(ctx context.Context, userIDFromToken, userIDReq uuid.UUID) error
	}

	projectService interface {
		CreateProject(ctx context.Context, projectReq *CreateProjectReq) (*model.Project, error)
		UpdateProject(ctx context.Context, projectReq *UpdateProjectReq) (*model.Project, error)
		DeleteProject(ctx context.Context, id int) error
		GetProjects(ctx context.Context, projectReq *GetProjectsReq) ([]model.Project, int, error)
		GetProjectInfo(ctx context.Context, id int) (*model.ProjectInfo, error)
		GetProjectCommits(ctx context.Context, id int) ([]model.CommitsInfo, error)
	}

	participantService interface {
		AddParticipant(ctx context.Context, isOwnerCreation bool, participant *AddedParticipant) (*model.Participant, error)
		UpdateParticipantRole(ctx context.Context, participant *ParticipantResp) (*model.Participant, error)
		DeleteParticipant(ctx context.Context, participantID int) error
		GetParticipantByID(ctx context.Context, id int) (*model.Participant, error)
		GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error)
		VerifyParticipant(ctx context.Context, userID uuid.UUID, projectID int) (*model.Participant, error)
		VerifyParticipantRole(ctx context.Context, userID uuid.UUID, projectID int, toAllow ...model.ParticipantRole) error
		VerifyParticipantByID(ctx context.Context, participantID int) (*model.Participant, error)
		VerifyParticipantRoleByID(ctx context.Context, participantID int, toAllow ...model.ParticipantRole) error
	}

	taskService interface {
		CreateTask(ctx context.Context, creatorUserID uuid.UUID, task *CreateTaskReq) (*model.Task, error)
		UpdateTask(ctx context.Context, taskReq *UpdateTaskReq) (*model.Task, error)
		DeleteTask(ctx context.Context, id int) error
		GetTasks(ctx context.Context, taskReq *GetTasksReq) ([]model.Task, int, error)
		GetTaskInfo(ctx context.Context, id int) (*model.TaskInfo, error)
	}

	OptionFunc func(s *Server)
)

func New(opts ...OptionFunc) *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := &Server{
		Server: &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  time.Duration(10) * time.Second,
			WriteTimeout: time.Duration(10) * time.Second},
	}
	for _, opt := range opts {
		opt(s)
	}

	rtr := gin.Default()
	cfg := cors.DefaultConfig()
	cfg.AllowHeaders = append(cfg.AllowHeaders, "Authorization", "Access-Control-Allow-Origin")
	cfg.AllowAllOrigins = true
	cors.New(cfg)
	rtr.Use(cors.New(cfg))
	// /api/*
	apiRtr := rtr.Group("/api")
	// /api/auth
	apiRtr.POST("/auth", s.auth)
	// /api/register
	apiRtr.POST("/register", s.register)

	// /api/user
	usersRtr := apiRtr.Group("/user")
	usersRtr.GET("/search", s.getPartialUsers)
	usersRtr.GET("/", s.authMiddleware(model.Admin, model.ProjectManager, model.Student), s.getUserProfileFromToken)
	usersRtr.GET("/:id", s.getUserProfile)
	usersRtr.PATCH("/", s.parseBodyToUpdatedUser, s.selfUpdateMiddleware(), s.updateUser)
	//usersRtr.DELETE("/:id", s.deleteUser)

	// /api/pm
	pmRtr := apiRtr.Group("/pm", s.authMiddleware(model.ProjectManager))
	pmRtr.POST("/", s.createProject)

	// /api/project
	projectRtr := apiRtr.Group("/project", s.authMiddleware(model.Admin, model.ProjectManager, model.Student))
	projectRtr.GET("/projects", s.getUserProjects)
	projectRtr.PATCH("/", s.parseBodyToUpdatedProject,
		s.verifyParticipantRoleMiddleware(model.RoleOwner, model.RoleTeamlead), s.updateProject)
	projectRtr.GET("/:projectId", s.getProjectInfo)
	projectRtr.GET("/:projectId/commits", s.getProjectCommits)
	projectRtr.GET("/:projectId/report", s.getProjectReport)
	projectRtr.DELETE("/remove", s.parseBodyToDeletedProject,
		s.verifyParticipantRoleMiddleware(model.RoleOwner), s.deleteProject)
	projectRtr.POST("/add-participant", s.parseBodyToAddedParticipant,
		s.verifyParticipantRoleMiddleware(model.RoleOwner, model.RoleTeamlead), s.addParticipant)
	projectRtr.PATCH("/update-participant", s.parseBodyToParticipantResp,
		s.verifyParticipantRoleMiddleware(model.RoleOwner), s.updateParticipant)
	projectRtr.DELETE("/remove-participant", s.parseBodyToParticipantResp,
		s.verifyParticipantRoleMiddleware(model.RoleOwner, model.RoleTeamlead), s.deleteParticipant)

	// /api/project/task
	taskRtr := projectRtr.Group("/:projectId/task", s.verifyParticipantMiddleware())
	taskRtr.POST("/", s.createTask)
	taskRtr.PATCH("/", s.updateTask)
	taskRtr.GET("/:taskId", s.getTaskInfo)
	taskRtr.DELETE("/", s.deleteTask)

	// /api/admin
	adminRtr := apiRtr.Group("/admin", s.authMiddleware(model.Admin))
	// /api/admin/users
	adminRtr.GET("/users/search", s.getFullUsers)
	adminRtr.GET("/users/search/:searchParam", s.getFullUsers)
	adminRtr.POST("/users", s.parseBodyToUpdatedUser, s.updateUser)
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
