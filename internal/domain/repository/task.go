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

func (r *Repository) GetTask(ctx context.Context, filter *TaskFilter) (*model.Task, error) {
	tasks, err := r.GetTasks(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	case len(tasks) == 0:
		return nil, ierr.ErrTaskNotFound
	default:
		return &tasks[0], nil
	}
}
func (r *Repository) GetTasks(ctx context.Context, filter *TaskFilter) ([]model.Task, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := sq.Select(
		"t.id", "t.name",
		"t.description", "t.suggested_estimate",
		"t.real_estimate", "t.participant_id",
		"t.creator_id", "t.status",
		"t.created_at", "t.updated_at",
		"t.project_id").
		From("tasks t").
		Where(conditionsFromTaskFilter(filter)).
		Limit(filter.Limit).
		Offset(filter.Offset).
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
	tasks := make([]model.Task, 0)
	for rows.Next() {
		task := model.Task{}
		if err = rows.Scan(
			&task.ID, &task.Name,
			&task.Description, &task.SuggestedEstimate,
			&task.RealEstimate, &task.ParticipantID,
			&task.CreatorID, &task.Status,
			&task.CreatedAt, &task.UpdatedAt,
			&task.ProjectID,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (r *Repository) GetTaskCountByFilter(ctx context.Context, filter *TaskFilter) (int, error) {
	var count int
	if err := r.sq.Select("COUNT(1)").
		From("tasks t").
		Where(conditionsFromTaskFilter(filter)).
		QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}
	return count, nil
}
func (r *Repository) InsertTask(ctx context.Context, task *model.Task) error {
	row := r.sq.Insert("tasks").
		Columns("name",
			"description", "suggested_estimate",
			"participant_id", "creator_id",
			"status", "project_id").
		Values(task.Name,
			task.Description, task.SuggestedEstimate,
			task.ParticipantID, task.CreatorID,
			task.Status, task.ProjectID).
		Suffix("RETURNING \"id\"").
		QueryRowContext(ctx)
	if err := row.Scan(&task.ID); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
	}
	return nil
}
func (r *Repository) UpdateTask(ctx context.Context, task *model.Task) error {
	_, err := r.sq.Update("tasks").
		SetMap(map[string]interface{}{
			"name":               task.Name,
			"description":        task.Description,
			"suggested_estimate": task.SuggestedEstimate,
			"real_estimate":      task.RealEstimate,
			"participant_id":     task.ParticipantID,
			"status":             task.Status,
			"updated_at":         task.UpdatedAt,
		}).Where(sq.Eq{"id": task.ID}).
		ExecContext(ctx)
	return err
}
func (r *Repository) DeleteTask(ctx context.Context, id int) error {
	_, err := r.sq.Delete("tasks").
		Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}

func (r *Repository) GetTaskInfo(ctx context.Context, id int) (*model.TaskInfo, error) {
	//получить фулл таску, получить юзера креатора и партисипанта (джоин через партисипантов и юзеров)
	r.sq.Select("t.id", "t.name", "t.description",
		"t.suggested_estimate", "t.real_estimate",
		"t.participant_id", "t.creator_id",
		"t.status", "t.created_at",
		"t.updated_at", "t.project_id",
		"u1.id", "u1.role",
		"u1.color_code", "u1.username",
		"u1.first_name", "u1.last_name",
		"u1.\"group\"", "u1.github_username",
		"u2.id", "u2.role",
		"u2.color_code", "u2.username",
		"u2.first_name", "u2.last_name",
		"u2.\"group\"", "u2.github_username").
		From("tasks t")
	// Join("participants part1 ON part1.id = t.creator_id").
	// Join("participants part2 ON u2.id = t.participant_id").
	return nil, nil
}
