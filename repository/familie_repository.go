package repository

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FamilieRepository interface {
	CreateNewFamilie(ctx context.Context, familie *model.Familie) error
	DeleteFamilie(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
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

func (r *repoFamilie) DeleteFamilie(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settings the transaction!")
	}

	query := `
		DELETE FROM families WHERE id = $1 AND created_by = $2;
		`

	if _, err := tx.ExecContext(ctx, query, id, user_id); err != nil {
		return errors.New("Failed to execute the query!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}
