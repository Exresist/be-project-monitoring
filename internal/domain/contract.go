package domain

import (
	"context"

	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/repository"

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
		// GetProjectInfo(ctx context.Context, id int, isTasks bool) (*model.ProjectInfo, error)
		GetProjectInfo(ctx context.Context, id int) (*model.ProjectInfo, error)
		GetProjectChecklist(ctx context.Context, id int) ([]model.Checklist, error)
		AddProjectChecklist(ctx context.Context, id int, checklist []model.Checklist) ([]model.Checklist, error)
		UpdateProjectChecklist(ctx context.Context, id int, checklist *model.Checklist) ([]model.Checklist, error)
		DeleteProjectChecklist(ctx context.Context, id int, itemID int) ([]model.Checklist, error)

		InsertProject(ctx context.Context, project *model.Project) error
		UpdateProject(ctx context.Context, project *model.Project) error
		DeleteProject(ctx context.Context, id int) error
	}

	participantRepo interface {
		AddParticipant(ctx context.Context, participant *model.Participant) error
		UpdateParticipantRole(ctx context.Context, participantID int, role string) error
		DeleteParticipant(ctx context.Context, id int) error
		GetParticipant(ctx context.Context, filter *repository.ParticipantFilter) (*model.Participant, error)
		GetParticipants(ctx context.Context, filter *repository.ParticipantFilter) ([]model.Participant, error)
	}

	taskRepo interface {
		GetTask(ctx context.Context, filter *repository.TaskFilter) (*model.Task, error)
		GetTasks(ctx context.Context, filter *repository.TaskFilter) ([]model.Task, error)
		GetCompletedTasksCountByGHUsername(ctx context.Context, projectID int) ([]model.TaskCount, error)
		GetTaskCountByFilter(ctx context.Context, filter *repository.TaskFilter) (int, error)
		GetTaskInfo(ctx context.Context, id int) (*model.TaskInfo, error)

		InsertTask(ctx context.Context, task *model.Task) error
		UpdateTask(ctx context.Context, task *model.Task) error
		DeleteTask(ctx context.Context, id int) error

		DeleteParticipantsFromTask(ctx context.Context, participantID int) error
	}
)
