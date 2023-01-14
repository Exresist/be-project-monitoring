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

// by ID escho nado
func (r *Repository) GetProject(ctx context.Context, filter *ProjectFilter) (*model.Project, error) {
	projects, err := r.GetProjects(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	case len(projects) == 0:
		return nil, ierr.ErrProjectNotFound
	default:
		return &projects[0], nil
	}
}

func (r *Repository) GetProjects(ctx context.Context, filter *ProjectFilter) ([]model.Project, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := sq.Select(
		"p.id", "p.name",
		"p.description", "p.photo_url",
		"p.report_url", "p.report_name",
		"p.repo_url", "p.active_to").
		From("projects p").
		Where(conditionsFromProjectFilter(filter)).
		Limit(filter.Limit).   // max = filter.Limit numbers
		Offset(filter.Offset). //  min = filter.Offset + 1
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)
	projects := make([]model.Project, 0)
	for rows.Next() {
		project := model.Project{}
		if err = rows.Scan(
			&project.ID, &project.Name,
			&project.Description, &project.PhotoURL,
			&project.ReportURL, &project.ReportName,
			&project.RepoURL, &project.ActiveTo,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *Repository) GetProjectCountByFilter(ctx context.Context, filter *ProjectFilter) (int, error) {
	var count int
	if err := r.sq.Select("COUNT(1)").
		From("projects p").
		Where(conditionsFromProjectFilter(filter)).
		QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}
	return count, nil
}

func (r *Repository) InsertProject(ctx context.Context, project *model.Project) error {
	row := r.sq.Insert("projects").
		Columns("name",
			"description", "photo_url",
			"active_to").
		Values(project.Name,
			project.Description, project.PhotoURL,
			project.ActiveTo).
		Suffix("RETURNING \"id\"").
		QueryRowContext(ctx)

	if err := row.Scan(&project.ID); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
	}
	return nil
}

func (r *Repository) UpdateProject(ctx context.Context, project *model.Project) error {
	_, err := r.sq.Update("projects").
		SetMap(map[string]interface{}{
			"name":        project.Name,
			"description": project.Description,
			"photo_url":   project.PhotoURL,
			"report_url":  project.ReportURL,
			"report_name": project.ReportName,
			"repo_url":    project.RepoURL,
			"active_to":   project.ActiveTo,
		}).ExecContext(ctx)
	return err
}
