package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AvraamMavridis/randomcolor"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"be-project-monitoring/internal/domain"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
)

type userStore struct {
	db        *sql.DB
	tableName string

	logger *zap.Logger
}

func NewUserStore(db *sql.DB, tableName string, logger *zap.Logger) domain.UserStore {
	return &userStore{
		db:        db,
		tableName: tableName,
		logger:    logger,
	}
}

func (u *userStore) GetByFilter(ctx context.Context, filter *domain.UserFilter) (*model.User, error) {
	users, err := u.GetListByFilter(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	case len(users) == 0:
		return nil, ierr.ErrUserNotFound
	default:
		return users[0], nil
	}
}

func (u *userStore) GetListByFilter(ctx context.Context, filter *domain.UserFilter) ([]*model.User, error) {
	rows, err := sq.Select(
		"id", "role",
		"color_code", "email",
		"username", "first_name",
		"last_name", "group",
		"github_username", "hashed_password").
		From(u.tableName).
		Where(u.conditions(filter)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		PlaceholderFormat(sq.Dollar).RunWith(u.db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func(res *sql.Rows) {
		err := res.Close()
		if err != nil {
			u.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}(rows)
	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Role,
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
func (u *userStore) GetCountByFilter(ctx context.Context, filter *domain.UserFilter) (int, error) {
	panic("")
}
func (u *userStore) DeleteByFilter(ctx context.Context, filter *domain.UserFilter) error {
	panic("")
}
func (u *userStore) Insert(ctx context.Context, user *model.User) (*model.User, error) {
	if err := sq.Insert("users").
		Columns("role", "color_code",
			"email", "username",
			"first_name", "last_name",
			"group", "github_username",
			"hashed_password").
		Values(user.Role, randomcolor.GetRandomColorInHex(),
			user.Email, user.Username,
			user.FirstName, user.LastName,
			user.Group, user.GithubUsername,
			user.HashedPassword).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).QueryRowContext(ctx).Scan(&user.ID); err != nil {
		return user, fmt.Errorf("error while inserting user: %w", err)
	}

	return user, nil
}
func (u *userStore) Update(ctx context.Context, user *model.User) error {
	panic("")
}

func (u *userStore) conditions(filter *domain.UserFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	if filter.IDs != nil {
		eq[u.tableName+".id"] = filter.IDs
	}
	if filter.Usernames != nil {
		eq[u.tableName+".username"] = filter.Usernames
	}
	if filter.Emails != nil {
		eq[u.tableName+".email"] = filter.Emails
	}

	return eq
}
