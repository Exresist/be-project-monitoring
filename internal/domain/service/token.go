package service

import (
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func (s *service) VerifyToken(ctx context.Context, token string, toAllow ...model.UserRole) error {
	// Parsing token fields
	claims, err := jwt.Parse(token, model.DecodeToken)
	if err != nil {
		return err
	}

	if !claims.Valid {
		return ierr.ErrInvalidToken
	}

	cl := claims.Claims.(jwt.MapClaims)
	roleID := cl["role"].(string)
	// Checking if role is in the list of the allowed roles
	for _, v := range toAllow {
		if roleID == string(v) {
			return nil
		}
	}
	return ierr.ErrAccessDenied
}

func (s *service) GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error) {
	// Parsing token fields
	claims, err := jwt.Parse(token, model.DecodeToken)
	if err != nil {
		return uuid.Nil, err
	}

	if !claims.Valid {
		return uuid.Nil, ierr.ErrInvalidToken
	}

	cl := claims.Claims.(jwt.MapClaims)
	userID := cl["id"].(uuid.UUID)

	return userID, nil
}

func (s *service) VerifySelf(ctx context.Context, token string, id uuid.UUID) error {
	tokenID, err := s.GetUserIDFromToken(ctx, token)
	if err != nil {
		return err
	}
	if tokenID != id {
		return ierr.ErrAccessDeniedAnotherUser
	}
	return nil
}
func (s *service) VerifyParticipant(ctx context.Context, userID uuid.UUID, projectID int) error {
	_, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(userID).ByProjectID(projectID))
	if err != nil {
		return err
	}
	return nil
}
func (s *service) VerifyParticipantRole(ctx context.Context, userID uuid.UUID, projectID int, toAllow ...model.ParticipantRole) error {
	participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(userID).ByProjectID(projectID))
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
