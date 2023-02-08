package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"

	"github.com/google/uuid"
)

func (s *service) AddParticipant(ctx context.Context, participantReq *api.AddParticipantReq) (*model.Participant, error) {
	if participantReq.ProjectID <= 0 {
		return nil, ierr.ErrInvalidProjectID
	}
	if participantReq.UserID == uuid.Nil {
		return nil, ierr.ErrInvalidUserID
	}
	if _, ok := model.ParticipantRoles[model.ParticipantRole(participantReq.Role)]; !ok {
		return nil, ierr.ErrInvalidParticipantRole
	}
	found, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(participantReq.UserID).ByProjectID(participantReq.ProjectID))
	if err != nil && err != ierr.ErrParticipantNotFound {
		return nil, err
	}
	if found != nil {
		return nil, ierr.ErrParticipantAlreadyExists
	}

	participant := &model.Participant{
		Role:      model.ParticipantRole(participantReq.Role),
		ProjectID: participantReq.ProjectID,
		User: model.User{
			ID: participantReq.UserID,
		},
	}
	return participant, s.repo.AddParticipant(ctx, participant)
}
func (s *service) GetParticipantByID(ctx context.Context, id int) (*model.Participant, error) {
	return s.repo.GetParticipant(ctx, repository.NewParticipantFilter().ByID(id))
}
func (s *service) GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error) {
	return s.repo.GetParticipants(ctx, repository.NewParticipantFilter().ByProjectID(projectID))
}

func (s *service) DeleteParticipant(ctx context.Context, userID uuid.UUID, projectID int) error {
	if projectID <= 0 {
		return ierr.ErrInvalidProjectID
	}
	if userID == uuid.Nil {
		return ierr.ErrInvalidUserID
	}
	participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(userID).ByProjectID(projectID))
	if err != nil{
		return err
	}
	if err := s.repo.DeleteParticipantsFromTask(ctx, participant.ID); err != nil {
		return err
	}
	return s.repo.DeleteParticipant(ctx, participant.ID)
}
