package repository

import (
	"be-project-monitoring/internal/db"

	"github.com/google/uuid"
)

type (
	UserFilter struct {
		IDs       []uuid.UUID `json:"id"`
		Usernames []string    `json:"username"`
		Emails    []string    `json:"email"`
		*db.Paginator
	}

	//TODO change me
	ProjectFilter struct {
		IDs          []int    `json:"id"`
		ProjectNames []string `json:"project_name"`
		*db.Paginator
	}
)

// UserFilter
func NewUserFilter() *UserFilter {
	return &UserFilter{Paginator: db.DefaultPaginator}
}
func (f *UserFilter) ByIDs(ids ...uuid.UUID) *UserFilter {
	f.IDs = ids
	return f
}
func (f *UserFilter) ByUsernames(usernames ...string) *UserFilter {
	if usernames != nil && usernames[0] != "" {
		f.Usernames = usernames
	}
	return f
}
func (f *UserFilter) ByEmails(emails ...string) *UserFilter {
	if emails != nil && emails[0] != "" {
		f.Emails = emails
	}
	return f
}
func (f *UserFilter) WithPaginator(limit, offset uint64) *UserFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

// ProjectFilter
func NewProjectFilter() *ProjectFilter {
	return &ProjectFilter{Paginator: db.DefaultPaginator}
}
func (f *ProjectFilter) ByIDs(ids ...int) *ProjectFilter {
	f.IDs = ids
	return f
}
func (f *ProjectFilter) ByProjectNames(projectNames ...string) *ProjectFilter {
	if projectNames != nil && projectNames[0] != "" {
		f.ProjectNames = projectNames
	}
	return f
}
func (f *ProjectFilter) WithPaginator(limit, offset uint64) *ProjectFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}
