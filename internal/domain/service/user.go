package service

import (
	"context"
	"errors"
	"fmt"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) CreateUser(ctx context.Context, userReq *api.CreateUserReq) (*model.User, string, error) {
	if userReq.Role == "" {
		userReq.Role = string(model.Student)
	}
	if _, ok := model.UserRoles[userReq.Role]; !ok {
		return nil, "", ierr.ErrInvalidUserRole
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
		ByEmail(user.Email).
		ByUsername(user.Username))
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

func (s *service) AuthUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByUsername(username))
	if err != nil {
		return "", fmt.Errorf("error while getting user: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", err
	}
	return model.GenerateToken(user)
}

func (s *service) GetUsers(ctx context.Context, userReq *api.GetUserReq) ([]model.User, int, error) {
	filter := repository.NewUserFilter().
		WithPaginator(uint64(userReq.Limit), uint64(userReq.Offset)).
		ByUsername(userReq.Username).ByEmail(userReq.Email)

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

func (s *service) UpdateUser(ctx context.Context, userReq *api.UpdateUserReq) (*model.User, error) {
	oldUser, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByID(userReq.ID))
	if err != nil {
		return nil, err
	}

	newUser, err := mergeUserFields(oldUser, userReq)
	if err != nil {
		return nil, err
	}
	return newUser, s.repo.UpdateUser(ctx, newUser)
}

func (s *service) DeleteUser(ctx context.Context, guid uuid.UUID) error {
	if _, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByID(guid)); err != nil {
		return err
	}
	return s.repo.DeleteUser(ctx, guid)
}

func (s *service) FindGithubUser(ctx context.Context, username string) bool {
	_, _, err := s.githubCl.Users.Get(ctx, username)
	return err == nil
}
func (s *service) GetUserProfile(ctx context.Context, guid uuid.UUID) (*model.Profile, error) {
	if _, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByID(guid)); err != nil {
		return nil, err
	}
	return s.repo.GetUserProfile(ctx, guid)
}

func hashPass(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

func mergeUserFields(oldUser *model.User, userReq *api.UpdateUserReq) (*model.User, error) {
	newUser := &model.User{
		ColorCode:      oldUser.ColorCode,
		Email:          oldUser.Email,
		ID:             userReq.ID,
		Role:           model.UserRole(*userReq.Role),
		Username:       *userReq.Username,
		FirstName:      *userReq.FirstName,
		LastName:       *userReq.LastName,
		Group:          *userReq.Group,
		GithubUsername: *userReq.GithubUsername,
		HashedPassword: hashPass(*userReq.Password),
	}

	if _, ok := model.UserRoles[*userReq.Role]; ok {
		newUser.Role = model.UserRole(*userReq.Role)
	} else {
		if userReq.Role == nil || *userReq.Role == "" {
			newUser.Role = oldUser.Role
		} else {
			return nil, ierr.ErrInvalidUserRole
		}
	}

	if userReq.Username == nil || *userReq.Username == "" {
		newUser.Username = oldUser.Username
	}
	if userReq.FirstName == nil || *userReq.FirstName == "" {
		newUser.FirstName = oldUser.FirstName
	}
	if userReq.LastName == nil || *userReq.LastName == "" {
		newUser.LastName = oldUser.LastName
	}
	if userReq.Group == nil || *userReq.Group == "" {
		newUser.Group = oldUser.Group
	}
	if userReq.GithubUsername == nil || *userReq.GithubUsername == "" {
		newUser.GithubUsername = oldUser.GithubUsername
	}
	if userReq.Password == nil || *userReq.Password == "" {
		newUser.HashedPassword = oldUser.HashedPassword
	}

	return newUser, nil
}
