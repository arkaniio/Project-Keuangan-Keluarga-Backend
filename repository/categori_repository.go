package repository

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateNewCategory(ctx context.Context, categories *model.Category) error
	UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type repoCategory struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repoCategory{db: db}
}

func (r *repoCategory) CreateNewCategory(ctx context.Context, categories *model.Category) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settings and add the transactions for this method!")
	}
	defer db_tx.Rollback()

	query := `
		INSERT INTO categories(id, user_id, name, type) 
		VALUES($1, $2, $3, $4);
	`

	rows, err := db_tx.ExecContext(ctx, query, categories.Id, categories.UserId, categories.Name, categories.Type)
	if err != nil {
		return errors.New("Failed to execute the db!" + err.Error())
	}

	last_infected, err := rows.RowsAffected()
	if err != nil {
		return errors.New("Failed to get the rows affected!")
	}

	if last_infected == 0 {
		return errors.New("Failed to insert the data!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the db!")
	}

	return nil

}

func (r *repoCategory) UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to add and settings the transactions")
	}
	defer db_tx.Rollback()

	full_query, args, err := utils.UpdateToolsCategory(payload, id)
	if err != nil {
		return errors.New("Failed to get the full query for categories!")
	}

	if _, err := db_tx.ExecContext(ctx, full_query, args...); err != nil {
		return errors.New("Failed to execute the db!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the db!")
	}

	return nil

}

func (r *repoCategory) DeleteCategory(ctx context.Context, id uuid.UUID) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to add and settings the transactions")
	}
	defer db_tx.Rollback()

	query := `
		DELETE FROM categories WHERE id = $1;
	`

	if _, err := db_tx.ExecContext(ctx, query, id); err != nil {
		return errors.New("Failed to execute the db!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the db!")
	}

	return nil

}
