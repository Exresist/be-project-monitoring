package repository

import (
	"be-project-monitoring/internal/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserFilter struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	*db.Paginator
}

func NewUserFilter() *UserFilter {
	return &UserFilter{Paginator: db.DefaultPaginator}
}

func (f *UserFilter) ByID(id uuid.UUID) *UserFilter {
	f.ID = id
	return f
}

func (f *UserFilter) ByUsername(username string) *UserFilter {
	f.Username = username
	return f
}

func (f *UserFilter) ByEmail(email string) *UserFilter {
	f.Email = email
	return f
}

func (f *UserFilter) WithPaginator(limit, offset uint64) *UserFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromUserFilter(filter *UserFilter) sq.Sqlizer {
	if filter.ID != uuid.Nil {
		return sq.Eq{"u.id": filter.ID}
	}
	usernameEq := make(sq.Eq)
	emailEq := make(sq.Eq)

	if filter.Username != "" {
		usernameEq["u.username"] = filter.Username
	}
	if filter.Email != "" {
		emailEq["u.email"] = filter.Email
	}

	return sq.Or{usernameEq, emailEq}

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
	ID   int
	Name string
	*db.Paginator
}

func NewProjectFilter() *ProjectFilter {
	return &ProjectFilter{Paginator: db.DefaultPaginator}
}

func (f *ProjectFilter) ByID(id int) *ProjectFilter {
	f.ID = id
	return f
}

func (f *ProjectFilter) ByProjectName(name string) *ProjectFilter {
	f.Name = name
	return f
}

func (f *ProjectFilter) WithPaginator(limit, offset uint64) *ProjectFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromProjectFilter(filter *ProjectFilter) sq.Sqlizer {
	if filter.ID > 0 {
		return sq.Eq{"p.id": filter.ID}
	}

	if filter.Name != "" {
		return sq.Eq{"p.name": filter.Name}
	}
	return nil
}

type TaskFilter struct {
	ID            int
	ProjectID     int
	ParticipantID *int
	Name          *string
	*db.Paginator
}

func NewTaskFilter() *TaskFilter {
	return &TaskFilter{Paginator: db.DefaultPaginator}
}

func (f *TaskFilter) ByID(id int) *TaskFilter {
	f.ID = id
	return f
}
func (f *TaskFilter) ByProjectID(id int) *TaskFilter {
	f.ProjectID = id
	return f
}
func (f *TaskFilter) ByParticipantID(id int) *TaskFilter {
	f.ParticipantID = &id
	return f
}
func (f *TaskFilter) ByTaskName(name string) *TaskFilter {
	f.Name = &name
	return f
}

func (f *TaskFilter) WithPaginator(limit, offset uint64) *TaskFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromTaskFilter(filter *TaskFilter) sq.Sqlizer {
	if filter.ID > 0 {
		return sq.Eq{"t.id": filter.ID}
	}

	projectEq := sq.Eq{"t.project_id": filter.ProjectID}
	nameEq := make(sq.Eq)
	participantEq := make(sq.Eq)
	if filter.Name != nil {
		nameEq["t.name"] = *filter.Name
	}
	if filter.ParticipantID != nil {
		participantEq["t.participant_id"] = *filter.ParticipantID
	}

	return sq.And{projectEq, participantEq, nameEq}
}

type ParticipantFilter struct {
	ID        int
	UserID    uuid.UUID
	ProjectID int
	*db.Paginator
}

func NewParticipantFilter() *ParticipantFilter {
	return &ParticipantFilter{Paginator: db.DefaultPaginator}
}
func (f *ParticipantFilter) ByID(id int) *ParticipantFilter {
	f.ID = id
	return f
}
func (f *ParticipantFilter) ByUserID(guid uuid.UUID) *ParticipantFilter {
	f.UserID = guid
	return f
}
func (f *ParticipantFilter) ByProjectID(id int) *ParticipantFilter {
	f.ProjectID = id
	return f
}
func (f *ParticipantFilter) WithPaginator(limit, offset uint64) *ParticipantFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}
func conditionsFromParticipantFilter(filter *ParticipantFilter) sq.Sqlizer {
	if filter.ID > 0 {
		return sq.Eq{"id": filter.ID}
	}
	if filter.UserID != uuid.Nil {
		return sq.Eq{"user_id": filter.UserID}
	}
	if filter.ProjectID > 0 {
		return sq.Eq{"project_id": filter.ProjectID}
	}
	return nil
}
