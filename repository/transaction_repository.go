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
	UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetTransactionById(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetAllTransaction(ctx context.Context) ([]model.PayloadTransactionWithCategory, error)
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
	defer db_tx.Rollback()

	query := `
		INSERT INTO transactions(id, user_id, type, amount, category_id, description, date, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	if transaction.Type != "expense" && transaction.Type != "income" {
		return errors.New("Failed to detect for a type, invalid type!")
	}

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

func (r *repoTransaction) DeleteTransaction(ctx context.Context, id uuid.UUID) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to get and settings the transactions!")
	}
	db_tx.Rollback()

	query := `
		DELETE FROM transactions WHERE id = $1;
	`

	if _, err := db_tx.ExecContext(ctx, query, id); err != nil {
		return errors.New("Failed to execute the query with context!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the db!")
	}

	return nil

}

func (r *repoTransaction) GetTransactionById(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {

	query := `
		SELECT id, user_id, type, amount, category_id, description, date, created_at, updated_at FROM transactions WHERE id = $1;
	`

	var transaction model.Transaction
	if err := r.db.GetContext(ctx, &transaction, query, id); err != nil {
		return nil, errors.New("Failed to get the transaction!")
	}

	return &transaction, nil

}

func (r *repoTransaction) GetAllTransaction(ctx context.Context) ([]model.PayloadTransactionWithCategory, error) {

	query := `
		SELECT t.id, t.user_id, t.type, t.amount, t.category_id, t.description, t.date, t.created_at, t.updated_at
		FROM transactions t
		JOIN categories c ON t.category_id = c.id;
	`

	var transaction_array []model.PayloadTransactionWithCategory
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, errors.New("Failed to load and query data transaction!")
	}

	for rows.Next() {
		var transaction_data model.PayloadTransactionDataCategory
		if err := rows.StructScan(transaction_data); err != nil {
			return nil, errors.New("Failed to get the transaction data and scan it into and struct")
		}
		second_rows, err := utils.PayloadJoinDataTransaction(transaction_data)
		if err != nil {
			return nil, errors.New("Failed to get the transaction data and join it into and struct")
		}
		transaction_array = append(transaction_array, *second_rows)
	}

	return transaction_array, nil

}
