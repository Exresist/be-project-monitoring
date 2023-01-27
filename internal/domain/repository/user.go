package repository

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"fmt"

	"github.com/AvraamMavridis/randomcolor"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *Repository) GetUser(ctx context.Context, filter *UserFilter) (*model.User, error) {
	users, err := r.GetUsers(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	case len(users) == 0:
		return nil, ierr.ErrUserNotFound
	default:
		return &users[0], nil
	}
}
func (r *Repository) GetUsers(ctx context.Context, filter *UserFilter) ([]model.User, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := r.sq.Select(
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "\"u.group\"",
		"u.github_username", "u.hashed_password").
		From("users u").
		Where(conditionsFromUserFilter(filter)).
		Limit(filter.Limit).   // max = filter.Limit numbers
		Offset(filter.Offset). //  min = filter.Offset + 1
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)
	users := make([]model.User, 0)
	for rows.Next() {
		user := model.User{}
		if err = rows.Scan(
			&user.ID, &user.Role,
			&user.ColorCode, &user.Email,
			&user.Username, &user.FirstName,
			&user.LastName, &user.Group,
			&user.GithubUsername, &user.HashedPassword,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *Repository) GetCountByFilter(ctx context.Context, filter *UserFilter) (int, error) {
	var count int
	if err := r.sq.Select("COUNT(1)").
		From("users u").
		Where(conditionsFromUserFilter(filter)).
		QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}
	return count, nil
}
func (r *Repository) InsertUser(ctx context.Context, user *model.User) error {
	_, err := r.sq.Insert("users").
		Columns("id", "role",
			"color_code", "email",
			"username", "first_name",
			"last_name", "\"group\"",
			"github_username", "hashed_password").
		Values(user.ID, user.Role,
			randomcolor.GetRandomColorInHex(), user.Email,
			user.Username, user.FirstName,
			user.LastName, user.Group,
			user.GithubUsername, user.HashedPassword).
		ExecContext(ctx)
	return err
}
func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.sq.Update("users").
		SetMap(map[string]interface{}{
			"role":            user.Role,
			"username":        user.Username,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"\"group\"":       user.Group,
			"github_username": user.GithubUsername,
			"hashed_password": user.HashedPassword,
		}).Where(sq.Eq{"id": user.ID}).
		ExecContext(ctx)
	return err
}
func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.sq.Delete("users").
		Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}
func (r *Repository) GetUserProfile(ctx context.Context, id uuid.UUID) (*api.GetUserProfileResp, error) {
	rows, err := r.sq.Select(
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "\"u.group\"",
		"u.github_username", "p.id",
		"p.name", "p.description",
		"p.active_to").
		From("participants part").
		Join("users u ON part.user_id = u.id").
		Join("projects p ON part.project_id = p.id").
		Where(sq.Eq{"id": id}).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)

	//ну тут кароче получится что на каждую запись для юзера с id=x в таблице participants (т.е. для каждого его проекта) userProject будет постоянно перезаписываться
	//можно сделать чтоб перед циклом 1 раз записался юзер в юзерПрофайл, а потом в цикле только проекты добавлялись
	userProfile := &api.GetUserProfileResp{}
	for rows.Next() {
		userProject := &api.UserProjectsResp{}
		if err = rows.Scan(
			&userProfile.ID, &userProfile.Role,
			&userProfile.ColorCode, &userProfile.Email,
			&userProfile.Username, &userProfile.FirstName,
			&userProfile.LastName, &userProfile.Group,
			&userProfile.GithubUsername, &userProject.ID,
			&userProject.Name, &userProject.Description,
			&userProject.ActiveTo,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		userProfile.UserProjects = append(userProfile.UserProjects, userProject)
	}
	return userProfile, nil
}
