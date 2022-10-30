package domain

import (
	"context"

	"github.com/google/uuid"

	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
)

type (
	UserStore interface {
		GetUser(ctx context.Context, filter *UserFilter) (*model.User, error)
		GetUsers(ctx context.Context, filter *UserFilter) ([]*model.User, error)
		GetCountByFilter(ctx context.Context, filter *UserFilter) (int, error)
		DeleteByFilter(ctx context.Context, filter *UserFilter) error

		Insert(ctx context.Context, user *model.User) error
		Update(ctx context.Context, user *model.User) error
	}

	UserFilter struct {
		IDs       []uuid.UUID `json:"id"`
		Usernames []string    `json:"username"`
		Emails    []string    `json:"email"`
		*db.Paginator
	}
)

func NewUserFilter() *UserFilter {
	return &UserFilter{Paginator: db.DefaultPaginator}
}

func (f *UserFilter) ByIDs(ids ...uuid.UUID) *UserFilter {
	f.IDs = ids
	return f
}

func (f *UserFilter) ByUsernames(usernames ...string) *UserFilter {
	f.Usernames = usernames
	return f
}

func (f *UserFilter) ByEmails(emails ...string) *UserFilter {
	f.Emails = emails
	return f
}

func (f *UserFilter) WithPaginator(limit, offset uint64) *UserFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}
