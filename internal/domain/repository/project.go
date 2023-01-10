package repository

import (
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type projectStore struct {
	db        *sql.DB
	tableName string

	logger *zap.SugaredLogger
}

func NewProjectStore(db *sql.DB, tableName string, logger *zap.SugaredLogger) *projectStore {
	return &projectStore{
		db:        db,
		tableName: tableName,
		logger:    logger,
	}
}

func (ps *projectStore) Insert(ctx context.Context, project *model.Project) error {
	_, err := sq.Insert("projects").
		Columns("id", "name",
			"description", "photo_url",
			"report_url", "report_name",
			"repo_url", "active_to").
		Values(project.ID, project.Name,
			project.Description, project.PhotoURL,
			project.ReportURL, project.ReportName,
			project.RepoURL, project.ActiveTo).
		PlaceholderFormat(sq.Dollar).
		RunWith(ps.db).ExecContext(ctx)
	return err
}

// by ID escho nado
func (ps *projectStore) GetProject(ctx context.Context, filter *ProjectFilter) (*model.Project, error) {
	users, err := ps.GetProjects(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	case len(users) == 0:
		return nil, ierr.ErrProjectNotFound
	default:
		return users[0], nil
	}
}

func (ps *projectStore) GetProjects(ctx context.Context, filter *ProjectFilter) ([]*model.Project, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := sq.Select(
		"id", "name",
		"description", "photo_url",
		"report_url", "report_name",
		"repo_url", "active_to").
		From(ps.tableName).
		Where(ps.conditions(filter)).
		Limit(filter.Limit).   // max = filter.Limit numbers
		Offset(filter.Offset). //  min = filter.Offset + 1
		PlaceholderFormat(sq.Dollar).RunWith(ps.db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			ps.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)
	projects := make([]*model.Project, 0)
	for rows.Next() {
		project := &model.Project{}
		err = rows.Scan(&project.ID, &project.Name,
			&project.Description, &project.PhotoURL,
			&project.ReportURL, &project.ReportName,
			&project.RepoURL, &project.ActiveTo)
		if err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (ps *projectStore) GetCountByFilter(ctx context.Context, filter *ProjectFilter) (int, error) {
	var count int
	if err := sq.Select("COUNT(1)").
		From(ps.tableName).
		Where(ps.conditions(filter)).
		PlaceholderFormat(sq.Dollar).
		RunWith(ps.db).QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}
	return count, nil
}

func (ps *projectStore) conditions(filter *ProjectFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	if filter.IDs != nil {
		eq[ps.tableName+".id"] = filter.IDs
	}
	// if len(filter.ProjectNames) != 0 && len(filter.Emails) != 0 {
	// 	usernameEq := make(sq.Eq)
	// 	emailEq := make(sq.Eq)
	// 	usernameEq[u.tableName+".username"] = filter.Usernames
	// 	emailEq[u.tableName+".email"] = filter.Emails
	// 	return sq.Or{eq, usernameEq, emailEq}
	// }
	if len(filter.ProjectNames) != 0 {
		eq[ps.tableName+".name"] = filter.ProjectNames
	}
	// if len(filter.Emails) != 0 {
	// 	eq[u.tableName+".email"] = filter.Emails
	// }

	return eq
}
