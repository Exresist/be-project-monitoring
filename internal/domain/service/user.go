package service

import (
	"context"

	"be-project-monitoring/internal/domain/model"
)

func (s *service) CreateUser(ctx context.Context, user *model.User) (*model.User, string, error) {
	user, err := s.store.Insert(ctx, user)
	if err != nil {
		return user, "", err
	}

	token, err := model.GenerateToken(user)
	return user, token, err
}
