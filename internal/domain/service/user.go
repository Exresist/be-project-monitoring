package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/repository"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
)

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, string, error) {
	found, err := s.userStore.GetUser(ctx, repository.NewUserFilter())
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

	if err = s.userStore.Insert(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := model.GenerateToken(user)
	return user, token, err
}

func (s *userService) AuthUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.userStore.GetUser(ctx, repository.NewUserFilter().ByUsernames(username))
	if err != nil {
		return "", fmt.Errorf("error while getting user: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", err
	}
	return model.GenerateToken(user)
}

func (s *userService) GetUsers(ctx context.Context, userReq *api.GetUserReq) ([]*model.User, int, error) {

	filter := repository.NewUserFilter().ByUsernames(userReq.Username).ByEmails(userReq.Email)
	filter.Limit = uint64(userReq.Limit)
	filter.Offset = uint64(userReq.Offset)

	count, err := s.userStore.GetCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.userStore.GetUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
