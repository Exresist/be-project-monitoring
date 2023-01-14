package repository

import (
	"database/sql"

	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db *sql.DB
	sq sq.StatementBuilderType

	logger *zap.SugaredLogger
}

func NewRepository(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db),
	}
}
