package repository

import (
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
}

type repoTransaction struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &repoTransaction{db: db}
}
