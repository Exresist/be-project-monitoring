package repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"be-project-monitoring/internal/domain/model"
)

func (r *Repository) GetProjectChecklist(ctx context.Context, id int) ([]model.Checklist, error) {
	rows, err := r.sq.Select(
		"id",
		"name",
		"project_id",
		"checked").
		From("checklist").
		Where(sq.Eq{"project_id": id}).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var checklist []model.Checklist

	for rows.Next() {
		item := model.Checklist{}
		if err = rows.Scan(&item.ID, &item.Name,
			&item.ProjectID, &item.Checked); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		checklist = append(checklist, item)
	}

	return checklist, nil
}

func (r *Repository) AddProjectChecklist(ctx context.Context, id int, checklist []model.Checklist) ([]model.Checklist, error) {
	q := r.sq.Insert("checklist").
		Columns("name",
			"project_id",
			"checked")

	for _, v := range checklist {
		q = q.Values(
			v.Name,
			id,
			v.Checked,
		)
	}

	if _, err := q.ExecContext(ctx); err != nil {
		return nil, err
	}

	return r.GetProjectChecklist(ctx, id)
}

func (r *Repository) UpdateProjectChecklist(ctx context.Context, id int, checklist *model.Checklist) ([]model.Checklist, error) {
	if _, err := r.sq.Update("checklist").
		Set("checked", checklist.Checked).
		Where(sq.Eq{"id": checklist.ID}).
		ExecContext(ctx); err != nil {
		return nil, err
	}

	return r.GetProjectChecklist(ctx, id)
}

func (r *Repository) DeleteProjectChecklist(ctx context.Context, id int, itemID int) ([]model.Checklist, error) {
	if _, err := r.sq.Delete("checklist").
		Where(sq.Eq{"id": itemID}).
		ExecContext(ctx); err != nil {
		return nil, err
	}

	return r.GetProjectChecklist(ctx, id)
}
