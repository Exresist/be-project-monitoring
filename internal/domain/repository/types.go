package repository

import (
	"be-project-monitoring/internal/db"

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

type ProjectFilter struct {
	Names []string
	*db.Paginator
}

func NewProjectFilter() *ProjectFilter {
	return &ProjectFilter{Paginator: db.DefaultPaginator}
}

func (f *ProjectFilter) ByProjectNames(names ...string) *ProjectFilter {
	f.Names = names
	return f
}

func (f *ProjectFilter) WithPaginator(limit, offset uint64) *ProjectFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromUserFilter(filter *UserFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	if filter.IDs != nil {
		eq["u.id"] = filter.IDs
	}
	if len(filter.Usernames) != 0 && len(filter.Emails) != 0 {
		usernameEq := make(sq.Eq)
		emailEq := make(sq.Eq)
		usernameEq["u.username"] = filter.Usernames
		emailEq["u.email"] = filter.Emails
		return sq.Or{eq, usernameEq, emailEq}
	}
	if filter.Usernames != nil {
		eq["u.username"] = filter.Usernames
	}
	if filter.Emails != nil {
		eq["u.email"] = filter.Emails
	}

	return eq
}

func conditionsFromProjectFilter(filter *ProjectFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	// if filter.IDs != nil {
	// 	eq[u.tableName+".id"] = filter.IDs
	// }
	// if len(filter.Usernames) != 0 && len(filter.Emails) != 0 {
	// 	usernameEq := make(sq.Eq)
	// 	emailEq := make(sq.Eq)
	// 	usernameEq[u.tableName+".username"] = filter.Usernames
	// 	emailEq[u.tableName+".email"] = filter.Emails
	// 	return sq.Or{eq, usernameEq, emailEq}
	// }
	// if len(filter.Usernames) != 0 {
	// 	eq[u.tableName+".username"] = filter.Usernames
	// }
	// if len(filter.Emails) != 0 {
	// 	eq[u.tableName+".email"] = filter.Emails
	// }

	return eq
}
