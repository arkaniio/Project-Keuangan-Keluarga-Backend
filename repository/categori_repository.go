package repository

import "github.com/jmoiron/sqlx"

type CategoryRepository interface {
}

type repoCategory struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repoCategory{db: db}
}
