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
	rows, err := r.sq.Select(
		"t.id", "t.name",
		"t.description", "t.suggested_estimate",
		"t.participant_id",
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
			&task.Description, &task.Estimate,
			&task.ParticipantID,
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
			task.Description, task.Estimate,
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
			"suggested_estimate": task.Estimate,
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
	rows, err := r.sq.Select("t.id", "t.name", "t.description",
		"t.suggested_estimate",
		"t.participant_id", "t.creator_id",
		"t.status", "t.created_at",
		"t.updated_at", "t.project_id",
		"u1.id", "u1.role",
		"u1.color_code", "u1.email",
		"u1.username", "u1.first_name",
		"u1.last_name", "u1.\"group\"",
		"u1.github_username",
		"u2.id", "u2.role",
		"u2.color_code", "u2.email",
		"u2.username", "u2.first_name",
		"u2.last_name", "u2.\"group\"",
		"u2.github_username").
		From("tasks t").
		LeftJoin("participants p_c ON p_c.id = t.creator_id").
		LeftJoin("participants p_p ON p_p.id = t.participant_id").
		LeftJoin("users u1 ON u1.id = p_c.user_id").
		LeftJoin("users u2 ON u2.id = p_p.user_id").
		Where(sq.Eq{"t.id": id}).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)

	if rows.Next() {
		taskInfo := model.TaskInfo{}
		if err := rows.Scan(&taskInfo.ID, &taskInfo.Name, &taskInfo.Description,
			&taskInfo.Estimate,
			&taskInfo.ParticipantID, &taskInfo.CreatorID,
			&taskInfo.Status, &taskInfo.CreatedAt,
			&taskInfo.UpdatedAt, &taskInfo.ProjectID,
			&taskInfo.Creator.ID, &taskInfo.Creator.Role,
			&taskInfo.Creator.ColorCode, &taskInfo.Creator.Email,
			&taskInfo.Creator.Username, &taskInfo.Creator.FirstName,
			&taskInfo.Creator.LastName, &taskInfo.Creator.Group,
			&taskInfo.Creator.GithubUsername,
			&taskInfo.Participant.ID, &taskInfo.Participant.Role,
			&taskInfo.Participant.ColorCode, &taskInfo.Participant.Email,
			&taskInfo.Participant.Username, &taskInfo.Participant.FirstName,
			&taskInfo.Participant.LastName, &taskInfo.Participant.Group,
			&taskInfo.Participant.GithubUsername,
		); err != nil {
			return nil, fmt.Errorf("error while performing sql request: %w", err)
		}
		return &taskInfo, nil
	}
	return nil, ierr.ErrTaskNotFound
}

func (r *Repository) DeleteParticipantsFromTask(ctx context.Context, participantID int) error {
	if _, err := r.sq.Update("tasks").
		SetMap(map[string]interface{}{
			"creator_id": "NULL",
		}).Where(sq.Eq{"creator_id": participantID}).
		ExecContext(ctx); err != nil {
		return err
	}
	_, err := r.sq.Update("tasks").
		SetMap(map[string]interface{}{
			"participant_id": "NULL",
		}).Where(sq.Eq{"participant_id": participantID}).
		ExecContext(ctx)
	return err
}
