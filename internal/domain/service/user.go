package service

import (
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) CreateUser(ctx context.Context, user *model.User) (*model.User, string, error) {
	found, err := s.repo.GetUser(ctx, repository.NewUserFilter().
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
		if found.GithubUsername == user.GithubUsername {
			return nil, "", ierr.ErrGithubUsernameAlreadyExists
		}
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, "", err
	}

	user.ID = userUUID

	if err = s.repo.Insert(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := model.GenerateToken(user)
	return user, token, err
}

func (s *service) AuthUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByUsernames(username))
	if err != nil {
		return "", fmt.Errorf("error while getting user: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", err
	}
	return model.GenerateToken(user)
}
