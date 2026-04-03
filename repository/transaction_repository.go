package repository

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	CreateNewTransactions(ctx context.Context, transactions *model.Transaction) error
}

type repoTransaction struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &repoTransaction{db: db}
}

func (r *repoTransaction) CreateNewTransactions(ctx context.Context, transaction *model.Transaction) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to setup the transactions settings!")
	}

	query := `
		INSERT INTO transactions(id, user_id, type, amount, category_id, description, date, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	rows, err := db_tx.ExecContext(ctx, query, transaction.Id, transaction.UserId, transaction.Type, transaction.Amount, transaction.CategoryId, transaction.Description, transaction.Date, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return errors.New("Failed to execute the db!")
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

func (r *repoTransaction) UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to get and settings the transactions!")
	}
	db_tx.Rollback()

	full_query, args, err := utils.UpdateToolsTransactions(payload, id)
	if err != nil {
		return errors.New("Failed to get the full query for this method!")
	}

	if _, err := db_tx.ExecContext(ctx, full_query, args...); err != nil {
		return errors.New("Failed to execute the query with context!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the db!")
	}

	return nil

}
