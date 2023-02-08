package repository

import (
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) AddParticipant(ctx context.Context, participant *model.Participant)  error {
	if _, err := r.sq.Insert("participants").
		Columns(
			"role",
			"user_id",
			"project_id",
		).
		Values(
			participant.Role,
			participant.User.ID,
			participant.ProjectID,
		).ExecContext(ctx); err != nil {
		return fmt.Errorf("error while saving participant: %w", err)
	}
	return nil
}

func (r *Repository) GetParticipant(ctx context.Context, filter *ParticipantFilter) (*model.Participant, error) {
	participants, err := r.GetParticipants(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get participant by id: %w", err)
	case len(participants) == 0:
		return nil, ierr.ErrParticipantNotFound
	default:
		return &participants[0], nil
	}
}
func (r *Repository) GetParticipants(ctx context.Context, filter *ParticipantFilter) ([]model.Participant, error) {
	rows, err := r.sq.Select(
		"p.id",
		"p.role",
		"p.user_id",
		"p.project_id",
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
		//Where("p.project_id = $1", projectID).QueryContext(ctx)
		Where(conditionsFromParticipantFilter(filter)).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while querying participants: %w", err)
	}
	participants := make([]model.Participant, 0)
	for rows.Next() {
		p := model.Participant{}
		if err := rows.Scan(
			&p.ID, &p.Role,
			&p.User.ID, &p.ProjectID,
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

func (r *Repository) DeleteParticipant(ctx context.Context, id int) error {	
	_, err := r.sq.Delete("participants").
		Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}

