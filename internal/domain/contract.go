package domain

import (
	"context"

	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
)

type (
	UserStore interface {
		GetByFilter(ctx context.Context, filter *UserFilter) (*model.User, error)
		GetListByFilter(ctx context.Context, filter *UserFilter) ([]*model.User, error)
		GetCountByFilter(ctx context.Context, filter *UserFilter) (int, error)
		DeleteByFilter(ctx context.Context, filter *UserFilter) error

		Insert(ctx context.Context, user *model.User) (*model.User, error)
		Update(ctx context.Context, user *model.User) error
	}

	UserFilter struct {
		IDs       []string `json:"id"`
		Usernames []string `json:"username"`
		Emails    []string `json:"email"`
		*db.Paginator
	}
)

func NewUserFilter() *UserFilter {
	return &UserFilter{Paginator: db.DefaultPaginator}
}

func (f *UserFilter) ByIDs(ids ...string) {
	f.IDs = ids
}

func (f *UserFilter) ByUsernames(usernames ...string) {
	f.Usernames = usernames
}

func (f *UserFilter) ByEmails(emails ...string) {
	f.Emails = emails
}

func (f *UserFilter) WithPaginator(limit, offset uint64) *UserFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}
