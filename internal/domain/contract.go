package domain

import (
	"context"

	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"

	"github.com/google/uuid"
)

type (
	Repository interface {
		userRepo
		projectRepo
		participantRepo
		taskRepo
	}

	userRepo interface {
		GetUser(ctx context.Context, filter *repository.UserFilter) (*model.User, error)
		GetFullUsers(ctx context.Context, filter *repository.UserFilter) ([]model.User, error)
		GetFullCountByFilter(ctx context.Context, filter *repository.UserFilter) (int, error)
		GetPartialUsers(ctx context.Context, filter *repository.UserFilter) ([]model.ShortUser, error)
		GetPartialCountByFilter(ctx context.Context, filter *repository.UserFilter) (int, error)
		GetUserProfile(ctx context.Context, id uuid.UUID) (*model.Profile, error)

		InsertUser(ctx context.Context, user *model.User) error
		UpdateUser(ctx context.Context, user *model.User) error
		DeleteUser(ctx context.Context, id uuid.UUID) error
	}

	projectRepo interface {
		GetProject(ctx context.Context, filter *repository.ProjectFilter) (*model.Project, error)
		GetProjects(ctx context.Context, filter *repository.ProjectFilter) ([]model.Project, error)
		GetProjectCountByFilter(ctx context.Context, filter *repository.ProjectFilter) (int, error)
		GetProjectInfo(ctx context.Context, id int) (*model.ProjectInfo, error)

		InsertProject(ctx context.Context, project *model.Project) error
		UpdateProject(ctx context.Context, project *model.Project) error
		DeleteProject(ctx context.Context, id int) error
	}

	participantRepo interface {
		AddParticipant(ctx context.Context, participant *model.Participant) error
		GetParticipant(ctx context.Context, filter *repository.ParticipantFilter) (*model.Participant, error)
		GetParticipants(ctx context.Context, filter *repository.ParticipantFilter) ([]model.Participant, error)
		DeleteParticipant(ctx context.Context, id int) error
	}

	taskRepo interface {
		GetTask(ctx context.Context, filter *repository.TaskFilter) (*model.Task, error)
		GetTasks(ctx context.Context, filter *repository.TaskFilter) ([]model.Task, error)
		GetTaskCountByFilter(ctx context.Context, filter *repository.TaskFilter) (int, error)
		GetTaskInfo(ctx context.Context, id int) (*model.TaskInfo, error)

		InsertTask(ctx context.Context, task *model.Task) error
		UpdateTask(ctx context.Context, task *model.Task) error
		DeleteTask(ctx context.Context, id int) error

		DeleteParticipantsFromTask(ctx context.Context, participantID int) error
	}
)
