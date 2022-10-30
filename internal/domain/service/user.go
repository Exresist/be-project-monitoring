package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"be-project-monitoring/internal/domain"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
)

func (s *service) CreateUser(ctx context.Context, user *model.User) (*model.User, string, error) {
	found, err := s.store.GetUser(ctx, domain.NewUserFilter().
		ByEmails(user.Email).
		ByUsernames(user.Username))
	if err != nil && !errors.Is(err, ierr.ErrUserNotFound) {
		return nil, "", err
	}
	if found != nil {
		if found.Email == user.Email {
			return nil, "", ierr.ErrEmailAlreadyExists
		}
		if found.Username == user.Username {
			return nil, "", ierr.ErrUsernameAlreadyExists
		}
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, "", err
	}

	found.ID = userUUID

	if err = s.store.Insert(ctx, found); err != nil {
		return nil, "", err
	}

	token, err := model.GenerateToken(found)
	return found, token, err
}

func (s *service) AuthUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.store.GetUser(ctx, domain.NewUserFilter().ByUsernames(username))
	if err != nil {
		return "", fmt.Errorf("error while getting user: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", err
	}
	return model.GenerateToken(user)
}
