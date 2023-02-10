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
	if _, ok := model.UserRoles[model.UserRole(userReq.Role)]; !ok {
		return nil, "", ierr.ErrInvalidUserRole
	}
	user := &model.User{
		ShortUser: model.ShortUser{
			Role:           model.UserRole(userReq.Role),
			Email:          userReq.Email,
			Username:       userReq.Username,
			FirstName:      userReq.FirstName,
			LastName:       userReq.LastName,
			Group:          userReq.Group,
			GithubUsername: userReq.GithubUsername,
		},
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

func (s *service) GetFullUsers(ctx context.Context, userReq *api.GetUserReq) ([]model.User, int, error) {
	filter := repository.NewUserFilter().
		WithPaginator(uint64(userReq.Limit), uint64(userReq.Offset)).
		ByUsername(userReq.Username).ByEmail(userReq.Email)

	count, err := s.repo.GetFullCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.repo.GetFullUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
func (s *service) GetPartialUsers(ctx context.Context, userReq *api.GetUserReq) ([]model.ShortUser, int, error) {
	if userReq.ProjectID <= 0 {
		return nil, 0, ierr.ErrInvalidProjectID
	}
	filter := repository.NewUserFilter().
		WithPaginator(uint64(userReq.Limit), uint64(userReq.Offset)).
		ByUsernameLike(userReq.Username).ByEmailLike(userReq.Email)
	if userReq.IsOnProject {
		filter = filter.ByAtProject(userReq.ProjectID)
	} else {
		filter = filter.ByNotAtProject(userReq.ProjectID)
	}

	count, err := s.repo.GetPartialCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	users, err := s.repo.GetPartialUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (s *service) UpdateUser(ctx context.Context, userReq *api.UpdateUserReq) (*model.User, error) {
	if userReq.ID == uuid.Nil {
		return nil, ierr.ErrInvalidUserID
	}
	oldUser, err := s.repo.GetUser(ctx, repository.NewUserFilter().ByID(userReq.ID))
	if err != nil {
		return nil, err
	}

	newUser, err := mergeUserFields(oldUser, userReq)
	if err != nil {
		return nil, err
	}
	if !s.FindGithubUser(ctx, newUser.GithubUsername) {
		return nil, ierr.ErrGithubUserNotFound
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
		ShortUser: model.ShortUser{
			ColorCode: oldUser.ColorCode,
			Email:     oldUser.Email,
			ID:        userReq.ID,
		},
	}

	if _, ok := model.UserRoles[model.UserRole(*userReq.Role)]; ok {
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
	} else {
		newUser.Username = *userReq.Username
	}
	if userReq.FirstName == nil || *userReq.FirstName == "" {
		newUser.FirstName = oldUser.FirstName
	} else {
		newUser.FirstName = *userReq.FirstName
	}
	if userReq.LastName == nil || *userReq.LastName == "" {
		newUser.LastName = oldUser.LastName
	} else {
		newUser.LastName = *userReq.LastName
	}
	if userReq.Group == nil || *userReq.Group == "" {
		newUser.Group = oldUser.Group
	} else {
		newUser.Group = *userReq.Group
	}
	if userReq.GithubUsername == nil || *userReq.GithubUsername == "" {
		newUser.GithubUsername = oldUser.GithubUsername
	} else {
		newUser.GithubUsername = *userReq.GithubUsername
	}
	if userReq.Password == nil || *userReq.Password == "" {
		newUser.HashedPassword = oldUser.HashedPassword
	} else {
		newUser.HashedPassword = hashPass(*userReq.Password)
	}

	return newUser, nil
}
