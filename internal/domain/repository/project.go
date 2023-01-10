package repository

import (
	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// by ID escho nado
func (r *Repository) GetProject(ctx context.Context, filter *ProjectFilter) (*model.Project, error) {
	users, err := r.GetProjects(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	case len(users) == 0:
		return nil, ierr.ErrProjectNotFound
	default:
		return users[0], nil
	}
}

func (r *Repository) GetProjects(ctx context.Context, filter *ProjectFilter) ([]*model.Project, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	// rows, err := sq.Select(
	// 	"id", "role",
	// 	"color_code", "email",
	// 	"username", "first_name",
	// 	"last_name", "\"group\"",
	// 	"github_username", "hashed_password").
	// 	From(r.tableName).
	// 	Where(r.conditions(filter)).
	// 	Limit(filter.Limit).   // max = filter.Limit numbers
	// 	Offset(filter.Offset). //  min = filter.Offset + 1
	// 	PlaceholderFormat(sq.Dollar).RunWith(r.db).QueryContext(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("error while performing sql request: %w", err)
	// }

	// defer func(res *sql.Rows) {
	// 	err = res.Close()
	// 	if err != nil {
	// 		r.logger.Error("error while closing sql rows", zap.Error(err))
	// 	}
	// }(rows)
	// users := make([]*model.Project, 0)
	// for rows.Next() {
	// 	user := &model.Project{}
	// 	err = rows.Scan(&user.ID, &user.Role,
	// 		&user.ColorCode, &user.Email,
	// 		&user.Username, &user.FirstName,
	// 		&user.LastName, &user.Group,
	// 		&user.GithubUsername, &user.HashedPassword)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error while scanning sql row: %w", err)
	// 	}
	// 	users = append(users, user)
	// }
	return nil, nil
}

func (r *Repository) GetProjectCountByFilter(ctx context.Context, filter *ProjectFilter) (int, error) {
	var count int
	// if err := sq.Select("COUNT(1)").
	// 	From(u.tableName).
	// 	Where(u.conditions(filter)).
	// 	PlaceholderFormat(sq.Dollar).
	// 	RunWith(u.db).QueryRowContext(ctx).Scan(&count); err != nil {
	// 	return 0, fmt.Errorf("error while scanning sql row: %w", err)
	// }
	return count, nil
}

func (r *Repository) conditions(filter *ProjectFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	// if filter.IDs != nil {
	// 	eq[u.tableName+".id"] = filter.IDs
	// }
	// if len(filter.Usernames) != 0 && len(filter.Emails) != 0 {
	// 	usernameEq := make(sq.Eq)
	// 	emailEq := make(sq.Eq)
	// 	usernameEq[u.tableName+".username"] = filter.Usernames
	// 	emailEq[u.tableName+".email"] = filter.Emails
	// 	return sq.Or{eq, usernameEq, emailEq}
	// }
	// if len(filter.Usernames) != 0 {
	// 	eq[u.tableName+".username"] = filter.Usernames
	// }
	// if len(filter.Emails) != 0 {
	// 	eq[u.tableName+".email"] = filter.Emails
	// }

	return eq
}
