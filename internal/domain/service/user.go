package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) CreateUser(ctx context.Context, userReq *api.CreateUserReq) (*model.User, string, error) {
	if _, ok := model.UserRoles[userReq.Role]; !ok {
		return nil, "", ierr.ErrInvalidRole
	}

	user := &model.User{
		Role:           model.UserRole(userReq.Role),
		Email:          userReq.Email,
		Username:       userReq.Username,
		FirstName:      userReq.FirstName,
		LastName:       userReq.LastName,
		Group:          userReq.Group,
		GithubUsername: userReq.GithubUsername,
		HashedPassword: hashPass(userReq.Password),
	}

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

	if err = s.repo.InsertUser(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := model.GenerateToken(user)
	return user, token, err
}

func (s *service) UpdateUser(ctx context.Context, userReq *api.UpdateUserReq) (*model.User, error) {
	oldUser, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByIDs(userReq.ID))
	if err != nil {
		return nil, fmt.Errorf("error while getting user: %w", err)
	}

	newUser := &model.User{
		ID:             userReq.ID,
		Role:           model.UserRole(userReq.Role),
		Username:       userReq.Username,
		FirstName:      userReq.FirstName,
		LastName:       userReq.LastName,
		Group:          userReq.Group,
		GithubUsername: userReq.GithubUsername,
		HashedPassword: hashPass(userReq.Password),
	}

	if err := mergeUserFields(oldUser, newUser); err != nil {
		return nil, err
	}

	return newUser, s.repo.UpdateUser(ctx, newUser)
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

func (s *service) GetUsers(ctx context.Context, userReq *api.GetUserReq) ([]model.User, int, error) {

	filter := repository.NewUserFilter().ByUsernames(userReq.Username).ByEmails(userReq.Email)
	filter.Limit = uint64(userReq.Limit)
	filter.Offset = uint64(userReq.Offset)

	count, err := s.repo.GetCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.repo.GetUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *service) FindGithubUser(ctx context.Context, username string) bool {
	_, _, err := s.githubCl.Users.Get(ctx, username)
	if err != nil {
		return false
	}
	return true
}

func hashPass(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

func mergeUserFields(oldUser, newUser *model.User) error {
	newUser.ColorCode = oldUser.ColorCode
	newUser.Email = oldUser.Email
	
	if newUser.Username == "" {
		newUser.Username = oldUser.Username
	}
	if newUser.FirstName == "" {
		newUser.FirstName = oldUser.FirstName
	}
	if newUser.LastName == "" {
		newUser.LastName = oldUser.LastName
	}
	if newUser.Group == "" {
		newUser.Group = oldUser.Group
	}
	if newUser.GithubUsername == "" {
		newUser.GithubUsername = oldUser.GithubUsername
	}
	if newUser.HashedPassword == "" {
		newUser.HashedPassword = oldUser.HashedPassword
	}
	if _, ok := model.UserRoles[string(newUser.Role)]; ok && newUser.Role == "" {
		newUser.Role = model.UserRole(newUser.Role)
	} else {
		return ierr.ErrInvalidRole
	}
	return nil
}
