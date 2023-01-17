package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
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
