package repository

import (
	"be-project-monitoring/internal/db"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserFilter struct {
	IDs       []uuid.UUID `json:"id"`
	Usernames []string    `json:"username"`
	Emails    []string    `json:"email"`
	*db.Paginator
}

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

func conditionsFromUserFilter(filter *UserFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	usernameEq := make(sq.Eq)
	emailEq := make(sq.Eq)
	if len(filter.IDs) != 0 {
		eq["u.id"] = filter.IDs
	}
	if len(filter.Usernames) != 0 {
		eq["u.username"] = filter.Usernames
	}
	if len(filter.Emails) != 0 {
		eq["u.email"] = filter.Emails
	}

	return sq.Or{eq, usernameEq, emailEq}

	// eq := make(sq.Eq)
	// if filter.IDs != nil {
	// 	eq["u.id"] = filter.IDs
	// }
	// if len(filter.Usernames) != 0 && len(filter.Emails) != 0 {
	// 	usernameEq := make(sq.Eq)
	// 	emailEq := make(sq.Eq)
	// 	usernameEq["u.username"] = filter.Usernames
	// 	emailEq["u.email"] = filter.Emails
	// 	return sq.Or{eq, usernameEq, emailEq}
	// }
	// if filter.Usernames != nil {
	// 	eq["u.username"] = filter.Usernames
	// }
	// if filter.Emails != nil {
	// 	eq["u.email"] = filter.Emails
	// }
	// return eq
}

type ProjectFilter struct {
	IDs   []int
	Names []string
	*db.Paginator
}

func NewProjectFilter() *ProjectFilter {
	return &ProjectFilter{Paginator: db.DefaultPaginator}
}

func (f *ProjectFilter) ByIDs(ids ...int) *ProjectFilter {
	f.IDs = ids
	return f
}

func (f *ProjectFilter) ByProjectNames(names ...string) *ProjectFilter {
	f.Names = names
	return f
}

func (f *ProjectFilter) WithPaginator(limit, offset uint64) *ProjectFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromProjectFilter(filter *ProjectFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	nameEq := make(sq.Eq)
	if len(filter.IDs) != 0 {
		eq["p.id"] = filter.IDs
	}
	if len(filter.Names) != 0 {
		eq["p.name"] = filter.Names
	}
	return sq.Or{eq, nameEq}
}

type TaskFilter struct {
	IDs   []int
	Names []string
	Dates []time.Time
	*db.Paginator
}

func NewTaskFilter() *TaskFilter {
	return &TaskFilter{Paginator: db.DefaultPaginator}
}

func (f *TaskFilter) ByIDs(ids ...int) *TaskFilter {
	f.IDs = ids
	return f
}

func (f *TaskFilter) ByTaskNames(names ...string) *TaskFilter {
	f.Names = names
	return f
}

func (f *TaskFilter) ByCreatedAt(dates ...time.Time) *TaskFilter {
	f.Dates = dates
	return f
}

func (f *TaskFilter) WithPaginator(limit, offset uint64) *TaskFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromTaskFilter(filter *TaskFilter) sq.Sqlizer {
	//TODO
	return nil
}
