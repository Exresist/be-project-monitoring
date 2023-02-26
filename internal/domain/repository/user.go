package repository

import (
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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
	// rows, err := r.sq.Select(
	// 	"u.id", "u.role",
	// 	"u.color_code", "u.email",
	// 	"u.username", "u.first_name",
	// 	"u.last_name", "u.\"group\"",
	// 	"u.github_username", "u.hashed_password").
	// 	From("users u").
	// 	Where(conditionsFromUserFilter(filter)).
	// 	Limit(filter.Limit).
	// 	Offset(filter.Offset).
	// 	QueryContext(ctx)

	// query := r.sq.Select(
	// 	"u.id", "u.role",
	// 	"u.color_code", "u.email",
	// 	"u.username", "u.first_name",
	// 	"u.last_name", "u.\"group\"",
	// 	"u.github_username", "u.hashed_password").
	// 	From("users u").
	// 	Where(conditionsFromUserFilter(filter))
	// fmt.Println(query.ToSql())
	rows, err := r.sq.Select(
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "u.\"group\"",
		"u.github_username", "u.hashed_password").
		From("users u").
		Where(conditionsFromUserFilter(filter)).
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
	fields := []string{
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "u.\"group\"",
		"u.github_username",
	}
	query := ""
	args := []interface{}{}
	err := error(nil)
	if filter.isOnProject {
		if query, args, err = getSQLStringForPartialUser(filter, fields...); err != nil {
			return nil, fmt.Errorf("error while generating sql query: %w", err)
		}
	} else {
		if query, args, err = getSQLStringForPartialUser(filter, fields[0]); err != nil {
			return nil, fmt.Errorf("error while generating sql query: %w", err)
		}

		query = `SELECT u.id, u.role, u.color_code, u.email,
		u.username, u.first_name, u.last_name, u."group", u.github_username
		FROM users u WHERE u.id NOT IN (` + query + `)`
	}
	if filter.likeText != "" {
		searchCondition, searchArgs, err := conditionsFromUserFilter(filter).ToSql()
		if err != nil {
			return nil, fmt.Errorf("error while adding conditions to sql row: %w", err)
		}
		searchCondition = strings.ReplaceAll(searchCondition, "?", "'"+searchArgs[0].(string)+"'")
		query = query + " AND " + searchCondition
	}
	//query = query + " LIMIT " + fmt.Sprint(filter.Limit) + " OFFSET " + fmt.Sprint(filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
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
	fields := []string{
		"u.id", "u.role",
		"u.color_code", "u.email",
		"u.username", "u.first_name",
		"u.last_name", "u.\"group\"",
		"u.github_username",
	}
	query := ""
	args := []interface{}{}
	err := error(nil)
	if filter.isOnProject {
		if query, args, err = getSQLStringForPartialUser(filter, "COUNT(1)"); err != nil {
			return 0, fmt.Errorf("error while generating sql query: %w", err)
		}
	} else {
		if query, args, err = getSQLStringForPartialUser(filter, fields[0]); err != nil {
			return 0, fmt.Errorf("error while generating sql query: %w", err)
		}
		query = `SELECT COUNT(1) FROM users u WHERE u.id NOT IN (` + query + `)`
	}

	if filter.likeText != "" {
		searchCondition, searchArgs, err := conditionsFromUserFilter(filter).ToSql()
		if err != nil {
			return 0, fmt.Errorf("error while adding conditions to sql row: %w", err)
		}
		searchCondition = strings.ReplaceAll(searchCondition, "?", "'"+searchArgs[0].(string)+"'")
		query = query + " AND " + searchCondition
	}
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}
	return count, nil
}
func getSQLStringForPartialUser(filter *UserFilter, fields ...string) (string, []interface{}, error) {
	return sq.Select(fields...).
		From("users u").
		Join("participants p on u.id = p.user_id").
		Where(conditionsFromUserFilterForProject(filter)).
		PlaceholderFormat(sq.Dollar).ToSql()
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
				ARRAY_AGG (p.photo_url) projects_photo_urls,
				ARRAY_AGG (p.active_to) projects_active_tos
			  FROM users u
			  LEFT JOIN participants part ON part.user_id = u.id
			  LEFT JOIN projects p ON part.project_id = p.id
			  WHERE u.id = $1
			  GROUP BY u.id, u.role, u.color_code, u.email, u.username,
					   u.first_name, u.last_name, u."group", u.github_username`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)

	if rows.Next() {
		profile := model.Profile{}
		projectsIDs := make(pq.ByteaArray, 0)
		projectsNames := make(pq.ByteaArray, 0)
		projectsDescriptions := make(pq.ByteaArray, 0)
		projectsPhotoURLs := make(pq.ByteaArray, 0)
		projectsActiveTos := make(pq.ByteaArray, 0)
		params := []any{&profile.ShortUser.ID, &profile.Role,
			&profile.ColorCode, &profile.Email,
			&profile.Username, &profile.FirstName,
			&profile.LastName, &profile.Group,
			&profile.GithubUsername, &projectsIDs,
			&projectsNames, &projectsDescriptions,
			&projectsPhotoURLs, &projectsActiveTos}

		if err := rows.Scan(params...); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		// for _, v := range params {
		// 	fmt.Printf("%v, %T \n\n", v, v)
		// }
		projects := make([]model.ShortProject, 0)
		for i := range projectsIDs {
			if projectsIDs[i] != nil {
				activeTo, err := time.Parse("2006-01-02", string(projectsActiveTos[i]))
				if err != nil {
					return nil, fmt.Errorf("error while parsing time: %w", err)
				}
				projectID, err := strconv.Atoi(string(projectsIDs[i]))
				if err != nil {
					return nil, err
				}

				shortProject := model.ShortProject{
					ID:       projectID,
					Name:     string(projectsNames[i]),
					ActiveTo: activeTo,
				}
				shortProject.Description.Scan(projectsDescriptions[i])
				shortProject.PhotoURL.Scan(projectsPhotoURLs[i])
				projects = append(projects, shortProject)
			}
		}
		profile.UserProjects = projects
		return &profile, nil
	}
	return nil, ierr.ErrUserNotFound
}
