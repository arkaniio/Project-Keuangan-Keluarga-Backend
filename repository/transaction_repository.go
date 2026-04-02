package repository

import (
	"github.com/jmoiron/sqlx"
)

type KeuanganRepository interface {
}

type repoKeuangan struct {
	db *sqlx.DB
}

func NewKeuanganRepository(db *sqlx.DB) KeuanganRepository {
	return &repoKeuangan{db: db}
}
