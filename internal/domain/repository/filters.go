package repository

import (
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserFilter struct {
	id              uuid.UUID
	username        string
	githubUsername  string
	email           string
	isProjectSearch bool
	isLike          bool
	likeText        string
	projectID       int
	isOnProject     bool
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
func (f *UserFilter) ByGithubUsername(githubUsername string) *UserFilter {
	f.githubUsername = githubUsername
	return f
}
func (f *UserFilter) ByLike(text string) *UserFilter {
	f.likeText = text
	f.isLike = true
	return f
}

//	func (f *UserFilter) ByUsernameLike(username string) *UserFilter {
//		f.username = username
//		f.isLike = true
//		return f
//	}
//
//	func (f *UserFilter) ByEmailLike(email string) *UserFilter {
//		f.email = email
//		f.isLike = true
//		return f
//	}

func (f *UserFilter) ByAtProject(projectID int) *UserFilter {
	f.projectID = projectID
	f.isOnProject = true
	f.isProjectSearch = true
	return f
}
func (f *UserFilter) ByNotAtProject(projectID int) *UserFilter {
	f.projectID = projectID
	f.isOnProject = false
	f.isProjectSearch = true
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
		githubEq := make(sq.Eq)
		if filter.username != "" && filter.email != "" && filter.githubUsername != "" {
			nameEq["u.username"] = filter.username
			emailEq["u.email"] = filter.email
			githubEq["u.github_username"] = filter.githubUsername
			return sq.Or{nameEq, emailEq, githubEq}
		}
		if filter.username != "" && filter.email != "" {
			nameEq["u.username"] = filter.username
			emailEq["u.email"] = filter.email
			return sq.Or{nameEq, emailEq}
		}
		if filter.username != "" && filter.githubUsername != "" {
			nameEq["u.username"] = filter.username
			githubEq["u.github_username"] = filter.githubUsername
			return sq.Or{nameEq, githubEq}
		}
		if filter.email != "" && filter.githubUsername != "" {
			emailEq["u.email"] = filter.email
			githubEq["u.github_username"] = filter.githubUsername
			return sq.Or{emailEq, githubEq}
		}
		if filter.username != "" {
			nameEq["u.username"] = filter.username
			return nameEq
		}
		if filter.email != "" {
			emailEq["u.email"] = filter.email
			return emailEq
		}
		if filter.githubUsername != "" {
			githubEq["u.github_username"] = filter.githubUsername
			return githubEq
		}
		return nil
	}
	or := make(sq.Or, 0)
	if filter.likeText != "" {
		or = sq.Or{sq.Like{"u.username": "%" + filter.likeText + "%"},
			sq.Like{"u.email": "%" + filter.likeText + "%"},
			sq.Like{"u.first_name": "%" + filter.likeText + "%"},
			sq.Like{"u.last_name": "%" + filter.likeText + "%"},
			sq.Like{"u.github_username": "%" + filter.likeText + "%"}}
	}
	if len(or) == 0 {
		return nil
	}
	return or

	// like := make(sq.Like)
	// if filter.username != "" {
	// 	like["u.username"] = "%" + filter.username + "%"
	// }
	// if filter.email != "" {
	// 	like["u.email"] = "%" + filter.email + "%"
	// }
	// if filter.isProjectSearch {
	// 	return sq.And{like, sq.Eq{"p.project_id": filter.projectID}}
	// }

	// return like
}
func conditionsFromUserFilterForProject(filter *UserFilter) sq.Sqlizer {
	return sq.Eq{"p.project_id": filter.projectID}
}

type ProjectFilter struct {
	ID     int
	Name   string
	isLike bool
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
func (f *ProjectFilter) ByProjectNameLike(name string) *ProjectFilter {
	f.Name = name
	f.isLike = true
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

	if !filter.isLike {
		if filter.Name != "" {
			return sq.Eq{"p.name": filter.Name}
		}
		return nil
	}
	return sq.Like{"p.name": "%" + filter.Name + "%"}
}

type TaskFilter struct {
	ID            int
	ProjectID     int
	ParticipantID *int
	Name          *string
	Status        model.TaskStatus
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

func (f *TaskFilter) ByStatus(status model.TaskStatus) *TaskFilter {
	f.Status = status
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
	Role      string
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
func (f *ParticipantFilter) ByRole(role string) *ParticipantFilter {
	f.Role = role
	return f
}
func (f *ParticipantFilter) WithPaginator(limit, offset uint64) *ParticipantFilter {
	f.Paginator = db.NewPaginator(limit, offset)
	return f
}
func conditionsFromParticipantFilter(filter *ParticipantFilter) sq.Sqlizer {
	if filter.ID > 0 && filter.ProjectID > 0 {
		return sq.Eq{"p.id": filter.ID,
			"p.project_id": filter.ProjectID}
	}
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
	if filter.Role != "" {
		eq["p.role"] = filter.Role
	}
	return eq
}
