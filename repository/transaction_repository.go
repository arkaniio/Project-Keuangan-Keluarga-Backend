package repository

import (
	"context"
	"errors"
	"fmt"
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
	GetAllTransaction(ctx context.Context, params model.PaginationParams) ([]model.PayloadTransactionWithCategory, int, error)
	GetAvgIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDay, error)
	GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDay, error)
	GetAvgIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeWeek, error)
	GetAvgExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseWeek, error)
	GetAvgIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeMonth, error)
	GetAvgExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseMonth, error)
	GetTransactionDataInExpenseType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error)
	GetTransactionDataInIncomeType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error)
	GetAvgExpenseDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDayNameCategory, error)
	GetAvgIncomeDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDayNameCategory, error)
	GetTotalExpenseByCategory(ctx context.Context, user_id uuid.UUID, category_id uuid.UUID) (int64, error)
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
		db_tx.Rollback()
		return errors.New("Failed to setup the transactions settings!")
	}

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

func (r *repoTransaction) GetTotalExpenseByCategory(ctx context.Context, userID, categoryID uuid.UUID) (int64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = $1
		AND category_id = $2
		AND type = 'expense'
	`
	var total int64
	err := r.db.GetContext(ctx, &total, query, userID, categoryID)
	return total, err
}

func (r *repoTransaction) UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		db_tx.Rollback()
		return errors.New("Failed to get and settings the transactions!")
	}

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

func (r *repoTransaction) GetAllTransaction(ctx context.Context, params model.PaginationParams) ([]model.PayloadTransactionWithCategory, int, error) {

	// ── Build dynamic WHERE clause ─────────────────────────────
	where := ""
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		where = fmt.Sprintf(" WHERE t.description ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// ── Count total items ──────────────────────────────────────
	countQuery := "SELECT COUNT(*) FROM transactions t JOIN categories c ON t.category_id = c.id" + where

	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, errors.New("Failed to count transactions: " + err.Error())
	}

	// ── Fetch paginated data ───────────────────────────────────
	offset := utils.CalculateOffset(params.Page, params.Limit)

	dataQuery := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.type, t.amount, t.category_id, t.description, t.date, t.created_at, t.updated_at
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		%s
		ORDER BY t.%s %s
		LIMIT $%d OFFSET $%d
	`, where, params.Sort, params.Order, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	var transaction_array []model.PayloadTransactionWithCategory
	rows, err := r.db.QueryxContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.New("Failed to load and query data transaction: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var transaction_data model.PayloadTransactionDataCategory
		if err := rows.StructScan(&transaction_data); err != nil {
			return nil, 0, errors.New("Failed to scan transaction data: " + err.Error())
		}
		second_rows, err := utils.PayloadJoinDataTransaction(transaction_data)
		if err != nil {
			return nil, 0, errors.New("Failed to parse transaction data: " + err.Error())
		}
		transaction_array = append(transaction_array, *second_rows)
	}

	return transaction_array, totalItems, nil

}

func (r *repoTransaction) GetAvgIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDay, error) {

	query := `
		SELECT DATE(date) as day,
			   AVG(amount) as avg_income
		FROM transactions
		WHERE user_id = $1 AND type = 'income'
		GROUP BY day
		ORDER BY day ASC;
	`

	var data model.AvgIncomeDay
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the income data svg!" + err.Error())
	}

	return &data, nil

}

func (r *repoTransaction) GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDay, error) {

	query := `
		SELECT DATE_TRUNC('day', date) as day,
			   AVG(amount) as avg_expense
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		GROUP BY day
		ORDER BY day ASC;
	`

	var data model.AvgExpenseDay
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("No data found!")
	}

	return &data, nil

}

func (r *repoTransaction) GetAvgIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeWeek, error) {

	query := `
		SELECT DATE_TRUNC('week', date) as week
			   AVG(amount) as avg_income_week
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		ORDER BY week
		GROUP BY week
	`

	var data model.AvgIncomeWeek
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg income week!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetAvgExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseWeek, error) {

	query := `
		SELECT DATE_TRUNC('week', date) as week
			   AVG(amount) as avg_expense_week
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		ORDER BY week
		GROUP BY week
	`

	var data model.AvgExpenseWeek
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg expense week!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetAvgIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeMonth, error) {

	query := `
		SELECT DATE_TRUNC('month', date) as month
			   AVG(amount) as avg_income_month
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		ORDER BY month
		GROUP BY month
	`

	var data model.AvgIncomeMonth
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg income month!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetAvgExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseMonth, error) {

	query := `
		SELECT DATE_TRUNC('month', date) as month
			   AVG(amount) as avg_expense_month
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		ORDER BY month
		GROUP BY month
	`

	var data model.AvgExpenseMonth
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg expense month!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetTransactionDataInExpenseType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {

	query := `
		SELECT id, user_id, type, amount, category_id, description, date, created_at, updated_at FROM transactions WHERE type = $1 AND user_id = $2;
	`

	var transaction model.Transaction
	if err := r.db.GetContext(ctx, &transaction, query, type_transaction); err != nil {
		return nil, errors.New("Failed to get the transaction!")
	}

	return &transaction, nil

}

func (r *repoTransaction) GetTransactionDataInIncomeType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {

	query := `
		SELECT id, user_id, type, amount, category_id, description, date, created_at, updated_at 
		FROM transactions WHERE type = $1;
	`

	var transaction model.Transaction
	if err := r.db.GetContext(ctx, &transaction, query, type_transaction); err != nil {
		return nil, errors.New("Failed to get the transaction!")
	}

	return &transaction, nil

}

func (r *repoTransaction) GetAvgExpenseDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDayNameCategory, error) {

	query := `
		SELECT DATE_TRUNC('day', t.date) as day 
		c.name as category
		AVG(t.amount) as avg_amount
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'expense'
		GROUP BY day
		ORDER BY day DESC
	`

	var data model.AvgExpenseDayNameCategory
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the data from db using query!")
	}

	return &data, nil

}

func (r *repoTransaction) GetAvgIncomeDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDayNameCategory, error) {

	query := `
		SELECT DATE_TRUNC('day', date) as day,
		c.name as category,
		AVG(t.amount) as avg_income
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'income'
		GROUP BY day, category
		ORDER BY day ASC;
	`

	var data model.AvgIncomeDayNameCategory
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the data from db using query!" + err.Error())
	}

	return &data, nil

}
