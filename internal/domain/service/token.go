package service

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
)

func (s *userService) VerifyToken(ctx context.Context, token string, toAllow ...model.UserRole) error {
	// Parsing token fields
	claims, err := jwt.Parse(token, model.DecodeToken)
	if err != nil {
		return err
	}

	if !claims.Valid {
		return ierr.ErrInvalidToken
	}

	cl := claims.Claims.(jwt.MapClaims)
	roleID := cl["role_id"].(model.UserRole)
	// Checking if role is in the list of the allowed roles
	for _, v := range toAllow {
		if roleID == v {
			return nil
		}
	}
	return ierr.ErrAccessDenied
}
