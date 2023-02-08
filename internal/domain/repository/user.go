package repository

import (
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AvraamMavridis/randomcolor"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func (r *Repository) GetUser(ctx context.Context, filter *UserFilter) (*model.User, error) {
	users, err := r.GetFullUsers(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	case len(users) == 0:
		return nil, ierr.ErrUserNotFound
	default:
		return &users[0], nil
	}
}
func (r *Repository) GetFullUsers(ctx context.Context, filter *UserFilter) ([]model.User, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := r.sq.Select(
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "u.\"group\"",
		"u.github_username", "u.hashed_password").
		From("users u").
		Where(conditionsFromUserFilter(filter)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
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
func (r *Repository) GetFullCountByFilter(ctx context.Context, filter *UserFilter) (int, error) {
	var count int
	if err := r.sq.Select("COUNT(1)").
		From("users u").
		Where(conditionsFromUserFilter(filter)).
		QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}
	return count, nil
}
func (r *Repository) GetPartialUsers(ctx context.Context, filter *UserFilter) ([]model.ShortUser, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := r.sq.Select(
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "u.\"group\"",
		"u.github_username").
		Distinct().
		From("users u").
		LeftJoin("participants p on p.user_id = u.id").
		Where(conditionsFromUserFilter(filter)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)
	users := make([]model.ShortUser, 0)
	for rows.Next() {
		user := model.ShortUser{}
		if err = rows.Scan(
			&user.ID, &user.Role,
			&user.ColorCode, &user.Email,
			&user.Username, &user.FirstName,
			&user.LastName, &user.Group,
			&user.GithubUsername,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *Repository) GetPartialCountByFilter(ctx context.Context, filter *UserFilter) (int, error) {
	var count int
	if err := r.sq.Select("COUNT(DISTINCT u.id)").
		Distinct().
		From("users u").
		LeftJoin("participants p on p.user_id = u.id").
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
func (r *Repository) GetUserProfile(ctx context.Context, id uuid.UUID) (*model.Profile, error) {
	query := `SELECT u.id, u.role, u.color_code, u.email, u.username,
	 			u.first_name, u.last_name, u."group", u.github_username,
				ARRAY_AGG (p.id) projects_ids, ARRAY_AGG (p.name) projects_names,
				ARRAY_AGG (p.description) projects_descriptions,
				ARRAY_AGG (p.active_to) projects_active_tos
			  FROM users u
			  JOIN participants part ON part.user_id = u.id
			  JOIN projects p ON part.project_id = p.id
			  WHERE u.id = $1
			  GROUP BY u.id, u.role, u.color_code, u.email, u.username,
					   u.first_name, u.last_name, u."group", u.github_username`

	projectsIDs := make(pq.Int64Array, 0)
	projectsNames := make(pq.StringArray, 0)
	projectsDescriptions := make(pq.StringArray, 0)
	projectsActiveTos := make([]time.Time, 0)

	profile := model.Profile{}
	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&profile.ID, &profile.Role,
		&profile.ColorCode, &profile.Email,
		&profile.Username, &profile.FirstName,
		&profile.LastName, &profile.Group,
		&profile.GithubUsername, &projectsIDs,
		&projectsNames, &projectsDescriptions,
		pq.Array(&projectsActiveTos)); err != nil {
		return nil, fmt.Errorf("error while scanning sql row: %w", err)
	}

	projects := make([]model.UserProjects, 0)
	for i := range projectsIDs {
		projects = append(projects, model.UserProjects{
			ID:          int(projectsIDs[i]),
			Name:        projectsNames[i],
			Description: projectsDescriptions[i],
			ActiveTo:    projectsActiveTos[i],
		})
	}
	profile.UserProjects = projects
	return &profile, nil
}
