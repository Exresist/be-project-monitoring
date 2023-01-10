package service

import (
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"

	"github.com/golang-jwt/jwt/v4"
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
