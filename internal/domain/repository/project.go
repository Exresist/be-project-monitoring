package repository

import (
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type projectStore struct {
	db        *sql.DB
	tableName string

	logger *zap.SugaredLogger
}

func NewProjectStore(db *sql.DB, tableName string, logger *zap.SugaredLogger) *projectStore {
	return &projectStore{
		db:        db,
		tableName: tableName,
		logger:    logger,
	}
}

func (u *projectStore) GetProject(ctx context.Context, filter *ProjectFilter) (*model.Project, error) {
	users, err := u.GetProjects(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	case len(users) == 0:
		return nil, ierr.ErrProjectNotFound
	default:
		return users[0], nil
	}
}

func (u *projectStore) GetProjects(ctx context.Context, filter *ProjectFilter) ([]*model.Project, error) {
	// rows, err := sq.Select(
	// 	"id", "role",
	// 	"color_code", "email",
	// 	"username", "first_name",
	// 	"last_name", "\"group\"",
	// 	"github_username", "hashed_password").
	// 	From(u.tableName).
	// 	Where(u.conditions(filter)).
	// 	Limit(filter.Limit).   // max = filter.Limit numbers
	// 	Offset(filter.Offset). //  min = filter.Offset + 1
	// 	PlaceholderFormat(sq.Dollar).RunWith(u.db).QueryContext(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("error while performing sql request: %w", err)
	// }

	// defer func(res *sql.Rows) {
	// 	err = res.Close()
	// 	if err != nil {
	// 		u.logger.Error("error while closing sql rows", zap.Error(err))
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

func (u *projectStore) conditions(filter *UserFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	if filter.IDs != nil {
		eq[u.tableName+".id"] = filter.IDs
	}
	if len(filter.Usernames) != 0 && len(filter.Emails) != 0 {
		usernameEq := make(sq.Eq)
		emailEq := make(sq.Eq)
		usernameEq[u.tableName+".username"] = filter.Usernames
		emailEq[u.tableName+".email"] = filter.Emails
		return sq.Or{eq, usernameEq, emailEq}
	}
	if len(filter.Usernames) != 0 {
		eq[u.tableName+".username"] = filter.Usernames
	}
	if len(filter.Emails) != 0 {
		eq[u.tableName+".email"] = filter.Emails
	}

	return eq
}