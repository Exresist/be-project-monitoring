package domain

import (
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	"context"
)

type (
	Repository interface {
		GetUser(ctx context.Context, filter *repository.UserFilter) (*model.User, error)
		GetUsers(ctx context.Context, filter *repository.UserFilter) ([]*model.User, error)
		GetCountByFilter(ctx context.Context, filter *repository.UserFilter) (int, error)
		DeleteByFilter(ctx context.Context, filter *repository.UserFilter) error

		Insert(ctx context.Context, user *model.User) error
		Update(ctx context.Context, user *model.User) error
	}
)
