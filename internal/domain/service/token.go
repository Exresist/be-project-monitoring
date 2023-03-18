package service

import (
	"context"

	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"

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
	roleID, ok := cl["role"].(string)
	if !ok {
		return ierr.ErrInvalidToken
	}
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
	return uuid.Parse(cl["id"].(string))
}
