package repository

import (
	"be-project-monitoring/internal/db"	

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserFilter struct {
	id          uuid.UUID
	username    string
	email       string
	isLike      bool
	projectID   int
	isOnProject bool
	*db.Paginator
}

func NewUserFilter() *UserFilter {
	return &UserFilter{Paginator: db.DefaultPaginator}
}

func (f *UserFilter) ByID(id uuid.UUID) *UserFilter {
	f.id = id
	return f
}

func (f *UserFilter) ByUsername(username string) *UserFilter {
	f.username = username
	return f
}
func (f *UserFilter) ByEmail(email string) *UserFilter {
	f.email = email
	return f
}
func (f *UserFilter) ByUsernameLike(username string) *UserFilter {
	f.username = username
	f.isLike = true
	return f
}
func (f *UserFilter) ByEmailLike(email string) *UserFilter {
	f.email = email
	f.isLike = true
	return f
}
func (f *UserFilter) ByAtProject(projectID int) *UserFilter {
	f.projectID = projectID
	f.isOnProject = true
	return f
}
func (f *UserFilter) ByNotAtProject(projectID int) *UserFilter {
	f.projectID = projectID
	f.isOnProject = false
	return f
}
func (f *UserFilter) WithPaginator(limit, offset uint64) *UserFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}

func conditionsFromUserFilter(filter *UserFilter) sq.Sqlizer {
	if filter.id != uuid.Nil {
		return sq.Eq{"u.id": filter.id}
	}
	if !filter.isLike {
		nameEq := make(sq.Eq)
		emailEq := make(sq.Eq)
		if filter.username != "" && filter.email != "" {
			nameEq["u.username"] = filter.username
			emailEq["u.email"] = filter.email
			return sq.Or{nameEq, emailEq}
		}
		if filter.username != "" {
			nameEq["u.username"] = filter.username
			return nameEq
		}
		if filter.email != "" {
			emailEq["u.email"] = filter.email
			return emailEq
		}
		return nil
	}

	like := make(sq.Like)
	if filter.username != "" {
		like["u.username"] = "%" + filter.username + "%"
	}
	if filter.email != "" {
		like["u.email"] = "%" + filter.email + "%"
	}
	if !filter.isOnProject && filter.projectID > 0 {
		return sq.And{like, sq.NotEq{"p.project_id": filter.projectID}} //проверить будет ли добавляться AND, если не будет like
	}
	if filter.isOnProject && filter.projectID > 0 {
		return sq.And{like, sq.Eq{"p.project_id": filter.projectID}} //проверить будет ли добавляться AND, если не будет like
	}
	return like
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

	// projectEq := sq.Eq{"t.project_id": filter.ProjectID}
	// nameEq := make(sq.Eq)
	// participantEq := make(sq.Eq)
	// if filter.Name != nil {
	// 	nameEq["t.name"] = *filter.Name
	// }
	// if filter.ParticipantID != nil {
	// 	participantEq["t.participant_id"] = *filter.ParticipantID
	// }

	// return sq.And{projectEq, participantEq, nameEq}

	//можно одну eq, т.к. под капотом используется AND, но проверить!
	eq := make(sq.Eq)
	eq["t.project_id"] = filter.ProjectID
	if filter.Name != nil {
		eq["t.name"] = *filter.Name
	}
	if filter.ParticipantID != nil {
		eq["t.participant_id"] = *filter.ParticipantID
	}
	return eq
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
		return sq.Eq{"p.id": filter.ID}
	}

	eq := make(sq.Eq) //можно одну eq, т.к. под капотом используется AND
	if filter.UserID != uuid.Nil {
		eq["p.user_id"] = filter.UserID
	}
	if filter.ProjectID > 0 {
		eq["p.project_id"] = filter.ProjectID
	}
	return eq
}
