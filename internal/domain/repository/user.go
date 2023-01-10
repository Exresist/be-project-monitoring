package repository

import (
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"fmt"

	"github.com/AvraamMavridis/randomcolor"
	sq "github.com/Masterminds/squirrel"
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
		return users[0], nil
	}
}

func (r *Repository) GetUsers(ctx context.Context, filter *UserFilter) ([]*model.User, error) {
	rows, err := r.sq.Select(
		"id", "role",
		"color_code", "email",
		"username", "first_name",
		"last_name", "\"group\"",
		"github_username", "hashed_password").
		From("users u").
		Where(conditionsFromUserFilter(filter)).
		Limit(filter.Limit).   // max = filter.Limit numbers
		Offset(filter.Offset). //  min = filter.Offset + 1
		PlaceholderFormat(sq.Dollar).RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)
	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.ID, &user.Role,
			&user.ColorCode, &user.Email,
			&user.Username, &user.FirstName,
			&user.LastName, &user.Group,
			&user.GithubUsername, &user.HashedPassword)
		if err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// TODO:
func (r *Repository) GetCountByFilter(ctx context.Context, filter *UserFilter) (int, error) {
	panic("TODO me")
}
func (r *Repository) DeleteByFilter(ctx context.Context, filter *UserFilter) error {
	panic("TODO me")
}
func (r *Repository) Insert(ctx context.Context, user *model.User) error {
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
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).ExecContext(ctx)
	return err
}
func (r *Repository) Update(ctx context.Context, user *model.User) error {
	panic("TODO me")
}
