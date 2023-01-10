package domain

import (
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	"context"
)

type (
	Repository interface {
		userRepo
		projectRepo
	}

	userRepo interface {
		GetUser(ctx context.Context, filter *repository.UserFilter) (*model.User, error)
		GetUsers(ctx context.Context, filter *repository.UserFilter) ([]*model.User, error)
		GetCountByFilter(ctx context.Context, filter *repository.UserFilter) (int, error)
		DeleteByFilter(ctx context.Context, filter *repository.UserFilter) error

		Insert(ctx context.Context, user *model.User) error
		Update(ctx context.Context, user *model.User) error
	}

	projectRepo interface {
		GetProject(ctx context.Context, filter *repository.ProjectFilter) (*model.Project, error)
		GetProjects(ctx context.Context, filter *repository.ProjectFilter) ([]*model.Project, error)
		GetProjectCountByFilter(ctx context.Context, filter *repository.ProjectFilter) (int, error)
	}
)
