package repository

import (
	"context"
	"fmt"

	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *Repository) AddParticipant(ctx context.Context, participant *model.Participant) error {
	if err := r.sq.Insert("participants").
		Columns(
			"role",
			"user_id",
			"project_id",
		).
		Values(
			participant.Role,
			participant.ShortUser.ID,
			participant.ProjectID,
		).
		Suffix("RETURNING \"id\"").
		QueryRowContext(ctx).
		Scan(&participant.ID); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
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
	).
		From("participants p").
		Join("users u ON u.id = p.user_id").
		Where(conditionsFromParticipantFilter(filter)).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while querying participants: %w", err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}()

	participants := make([]model.Participant, 0)
	for rows.Next() {
		p := model.Participant{}
		if err = rows.Scan(
			&p.ID, &p.Role,
			&p.ShortUser.ID, &p.ProjectID,
			&p.ShortUser.Role, &p.ColorCode,
			&p.Email, &p.Username,
			&p.FirstName, &p.LastName,
			&p.Group, &p.GithubUsername,
		); err != nil {
			return nil, fmt.Errorf("error while scanning row: %w", err)
		}

		participants = append(participants, p)
	}
	return participants, nil
}

func (r *Repository) UpdateParticipantRole(ctx context.Context, participantID int, role string) error {
	_, err := r.sq.Update("participants").
		SetMap(map[string]interface{}{
			"role": role,
		}).Where(sq.Eq{"id": participantID}).
		ExecContext(ctx)
	return err
}

func (r *Repository) DeleteParticipant(ctx context.Context, id int) error {
	_, err := r.sq.Delete("participants").
		Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}
