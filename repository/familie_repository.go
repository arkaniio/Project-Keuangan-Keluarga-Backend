package repository

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"

	"github.com/jmoiron/sqlx"
)

type FamilieRepository interface {
	CreateNewFamilie(ctx context.Context, familie *model.Familie) error
}

type repoFamilie struct {
	db *sqlx.DB
}

func NewFamilieRepository(db *sqlx.DB) FamilieRepository {
	return &repoFamilie{db: db}
}

func (r *repoFamilie) CreateNewFamilie(ctx context.Context, familie *model.Familie) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to adding the new transaction for this method!")
	}

	query := `
		INSERT INTO families(id, name, created_by, created_at)
		VALUES ($1, $2, $3, $4);
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return errors.New("Failed to execute the context!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}
