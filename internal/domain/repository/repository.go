package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db *sql.DB
	sq sq.StatementBuilderType
}

func NewRewardsRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db),
	}
}
