package repository

import (
	"be-project-monitoring/internal/domain/model"
	"context"
	"fmt"
)

func (r *Repository) AddParticipant(ctx context.Context, participant *model.Participant) ([]model.Participant, error) {
	if _, err := r.sq.Insert("participants").
		Columns(
			"role",
			"user_id",
			"project_id",
		).
		Values(
			participant.Role,
			participant.UserID,
			participant.ProjectID,
		).ExecContext(ctx); err != nil {
		return nil, fmt.Errorf("error while saving participant: %w", err)
	}

	return r.GetParticipants(ctx, participant.ProjectID)
}

func (r *Repository) GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error) {
	rows, err := r.sq.Select(
		"p.role",
		"p.user_id",
		"p.project_id",
		"u.id",
		"u.role",
		"u.color_code",
		"u.email",
		"u.username",
		"u.first_name",
		"u.last_name",
		"u.group",
		"u.github_username",
		"u.hashed_password",
	).
		From("participants p").
		Join("users u ON u.id = p.user_id").
		Where("p.project_id = $1", projectID).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while querying participants: %w", err)
	}
	participants := make([]model.Participant, 0)
	for rows.Next() {
		p := model.Participant{}
		if err := rows.Scan(
			&p.Role, &p.UserID,
			&p.ProjectID, &p.ID,
			&p.User.Role, &p.ColorCode,
			&p.Email, &p.Username,
			&p.FirstName, &p.LastName,
			&p.Group, &p.GithubUsername,
			&p.HashedPassword,
		); err != nil {
			return nil, fmt.Errorf("error while scanning row: %w", err)
		}
		participants = append(participants, p)
	}
	return participants, nil
}
