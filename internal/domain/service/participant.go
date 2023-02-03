package service

import (
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"

	"github.com/google/uuid"
)

func (s *service) AddParticipant(ctx context.Context, participant *model.Participant) ([]model.Participant, error) {
	if participant.ProjectID <= 0 {
		return nil, ierr.ErrInvalidProjectID
	}
	if participant.User.ID == uuid.Nil {
		return nil, ierr.ErrInvalidUserID
	}
	if _, ok := model.ParticipantRoles[participant.Role]; !ok {
		return nil, ierr.ErrInvalidParticipantRole
	}
	return s.repo.AddParticipant(ctx, participant)
}
func (s *service) GetParticipantByID(ctx context.Context, id int) (*model.Participant, error) {
	return s.repo.GetParticipant(ctx, repository.NewParticipantFilter().ByID(id))
}
func (s *service) GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error) {
	return s.repo.GetParticipants(ctx, repository.NewParticipantFilter().ByProjectID(projectID))
}
