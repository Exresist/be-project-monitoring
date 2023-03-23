package service

import (
	"context"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"be-project-monitoring/internal/repository"

	"github.com/google/uuid"
)

func (s *service) AddParticipant(ctx context.Context, isOwnerCreation bool, participantReq *api.AddedParticipant) (*model.Participant, error) {
	if participantReq.ProjectID <= 0 {
		return nil, ierr.ErrInvalidProjectID
	}

	if participantReq.UserID == uuid.Nil {
		return nil, ierr.ErrInvalidUserID
	}

	if _, ok := model.ParticipantRoles[model.ParticipantRole(participantReq.Role)]; !ok ||
		!isOwnerCreation && participantReq.Role == string(model.RoleOwner) {
		return nil, ierr.ErrInvalidParticipantRole
	}

	var (
		isUser     bool
		teamLeadID int
	)

	if participants, err := s.repo.GetParticipants(ctx, repository.NewParticipantFilter().
		ByProjectID(participantReq.ProjectID)); err != nil {
		return nil, err
	} else if len(participants) != 0 {
		for _, v := range participants {

			if v.ShortUser.ID == participantReq.UserID {
				isUser = true
			}

			if v.Role == model.RoleTeamlead {
				teamLeadID = v.ID
			}
		}
	}

	if isUser {
		return nil, ierr.ErrParticipantAlreadyExists
	}

	if participantReq.Role == string(model.RoleTeamlead) && teamLeadID != 0 {
		if err := s.repo.UpdateParticipantRole(ctx, teamLeadID, string(model.RoleParticipant)); err != nil {
			return nil, err
		}
	}

	participant := &model.Participant{
		Role:      model.ParticipantRole(participantReq.Role),
		ProjectID: participantReq.ProjectID,
		ShortUser: model.ShortUser{
			ID: participantReq.UserID,
		},
	}

	if err := s.repo.AddParticipant(ctx, participant); err != nil {
		return nil, err
	}

	return s.GetParticipantByID(ctx, participant.ID)
}

func (s *service) UpdateParticipantRole(ctx context.Context, participant *api.ParticipantResp) (*model.Participant, error) {

	if _, ok := model.ParticipantRoles[model.ParticipantRole(participant.Role)]; !ok ||
		participant.Role == string(model.RoleOwner) {
		return nil, ierr.ErrInvalidParticipantRole
	}

	var (
		isParticipant bool
		teamLeadID    int
	)

	if participants, err := s.repo.GetParticipants(ctx, repository.NewParticipantFilter().
		ByProjectID(participant.ProjectID)); err != nil {
		return nil, err
	} else if len(participants) == 0 {
		return nil, ierr.ErrParticipantsNotFound
	} else {
		for _, v := range participants {
			if v.ID == participant.ID {
				isParticipant = true
			}
			if v.Role == model.RoleTeamlead {
				teamLeadID = v.ID
			}
		}
	}

	if !isParticipant {
		return nil, ierr.ErrParticipantNotFound
	}

	if participant.Role == string(model.RoleTeamlead) && teamLeadID != 0 {
		if err := s.repo.UpdateParticipantRole(ctx, teamLeadID, string(model.RoleParticipant)); err != nil {
			return nil, err
		}
	}

	return &model.Participant{
		ID:        participant.ID,
		Role:      model.ParticipantRole(participant.Role),
		ProjectID: participant.ProjectID,
		ShortUser: participant.User,
	}, s.repo.UpdateParticipantRole(ctx, participant.ID, participant.Role)
}

func (s *service) GetParticipantByID(ctx context.Context, id int) (*model.Participant, error) {
	return s.repo.GetParticipant(ctx, repository.NewParticipantFilter().ByID(id))
}

func (s *service) GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error) {
	return s.repo.GetParticipants(ctx, repository.NewParticipantFilter().ByProjectID(projectID))
}

func (s *service) DeleteParticipant(ctx context.Context, participantID int) error {

	if err := s.repo.DeleteParticipantsFromTask(ctx, participantID); err != nil {
		return err
	}

	return s.repo.DeleteParticipant(ctx, participantID)
}

func (s *service) VerifyParticipant(ctx context.Context, userID uuid.UUID, projectID int) (*model.Participant, error) {

	participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(userID).ByProjectID(projectID))
	if err != nil {
		return nil, ierr.ErrUserIsNotOnProject
	}

	return participant, nil
}
func (s *service) VerifyParticipantRole(ctx context.Context, userID uuid.UUID, projectID int, toAllow ...model.ParticipantRole) error {

	participant, err := s.VerifyParticipant(ctx, userID, projectID)
	if err != nil {
		return err
	}

	// Checking if role is in the list of the allowed roles
	for _, v := range toAllow {
		if participant.Role == v {
			return nil
		}
	}

	return ierr.ErrAccessDeniedWrongParticipantRole
}
func (s *service) VerifyParticipantByID(ctx context.Context, participantID int) (*model.Participant, error) {

	participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByID(participantID))
	if err != nil {
		return nil, ierr.ErrUserIsNotOnProject
	}

	return participant, nil
}
func (s *service) VerifyParticipantRoleByID(ctx context.Context, participantID int, toAllow ...model.ParticipantRole) error {

	participant, err := s.VerifyParticipantByID(ctx, participantID)
	if err != nil {
		return err
	}

	for _, v := range toAllow {
		if participant.Role == v {
			return nil
		}
	}

	return ierr.ErrAccessDeniedWrongParticipantRole
}
