package repository

import (
	"context"
	"database/sql"
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
	GetTransactionByUserId(ctx context.Context, user_id uuid.UUID) (*model.Transaction, error)
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
	GetTotalExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseDay, error)
	GetTotalExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseWeek, error)
	GetTotalExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseMonth, error)
	GetTotalIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeDay, error)
	GetTotalIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeWeek, error)
	GetTotalIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeMonth, error)
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
		INSERT INTO transactions(id, user_id, family_member_id, type, amount, category_id, description, date, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	if transaction.Type != "expense" && transaction.Type != "income" {
		return errors.New("Failed to detect for a type, invalid type!")
	}

	rows, err := db_tx.ExecContext(ctx, query, transaction.Id, transaction.UserId, transaction.FamilyMemberId, transaction.Type, transaction.Amount, transaction.CategoryId, transaction.Description, transaction.Date, transaction.CreatedAt, transaction.UpdatedAt)
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

func (r *repoTransaction) GetTransactionByUserId(ctx context.Context, user_id uuid.UUID) (*model.Transaction, error) {

	query := `
		SELECT id, user_id, type, amount, category_id, description, date, created_at, updated_at
		FROM transactions WHERE user_id = $1;
	`

	var data model.Transaction
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the transaction!")
	}

	return &data, nil

}

func (r *repoTransaction) DeleteTransaction(ctx context.Context, id uuid.UUID) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to get and settings the transactions!")
	}
	defer db_tx.Rollback()

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
		SELECT t.id, t.user_id, t.family_member_id, t.type, t.amount, t.category_id, t.description, t.date, t.created_at, t.updated_at, c.name, c.type as category_type
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
		SELECT DATE_TRUNC('week', date) as week,
			   AVG(amount) as income
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		GROUP BY week
		ORDER BY week
	`

	var data model.AvgIncomeWeek
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg income week!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetAvgExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseWeek, error) {

	query := `
		SELECT DATE_TRUNC('week', date) as week,
			   AVG(amount) as expense
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		GROUP BY week
		ORDER BY week
	`

	var data model.AvgExpenseWeek
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg expense week!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetAvgIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeMonth, error) {

	query := `
		SELECT DATE_TRUNC('month', date) as month,
			   AVG(amount) as income
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		GROUP BY month
		ORDER BY month
	`

	var data model.AvgIncomeMonth
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the avg income month!" + err.Error())
	}

	return &data, nil
}

func (r *repoTransaction) GetAvgExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseMonth, error) {

	query := `
		SELECT DATE_TRUNC('month', date) as month,
			   AVG(amount) as expense
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		GROUP BY month
		ORDER BY month
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
	if err := r.db.GetContext(ctx, &transaction, query, type_transaction, user_id); err != nil {
		return nil, errors.New("Failed to get the transaction!")
	}

	return &transaction, nil

}

func (r *repoTransaction) GetTransactionDataInIncomeType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {

	query := `
		SELECT id, user_id, type, amount, category_id, description, date, created_at, updated_at 
		FROM transactions WHERE type = $1 AND user_id = $2;
	`

	var transaction model.Transaction
	if err := r.db.GetContext(ctx, &transaction, query, type_transaction, user_id); err != nil {
		return nil, errors.New("Failed to get the transaction!")
	}

	return &transaction, nil

}

func (r *repoTransaction) GetAvgExpenseDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDayNameCategory, error) {

	query := `
		SELECT DATE_TRUNC('day', t.date) as day,
		c.name as name,
		AVG(t.amount) as avg_expense
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'expense'
		GROUP BY day, c.name
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

func (r *repoTransaction) GetTotalExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseDay, error) {

	query := `
		SELECT DATE_TRUNC('day', date) as day,
		SUM(amount) as total_expense
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		GROUP BY day
		ORDER BY day ASC;
	`

	var data model.TotalExpenseDay
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No rows detected!")
		}
		return nil, errors.New("Failed to get the total expense day!" + err.Error())
	}

	return &data, nil

}

func (r *repoTransaction) GetTotalExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseWeek, error) {

	query := `
		SELECT DATE_TRUNC('week', date) as week,
		SUM(amount) as total_expense
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		GROUP BY week
		ORDER BY week ASC;
	`

	var data model.TotalExpenseWeek
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No rows detected!")
		}
		return nil, errors.New("Failed to get the total expense week!" + err.Error())
	}

	return &data, nil

}

func (r *repoTransaction) GetTotalExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseMonth, error) {

	query := `
		SELECT DATE_TRUNC('month', date) as month,
		SUM(amount) as total_expense
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense'
		GROUP BY month
		ORDER BY month ASC;
	`

	var data model.TotalExpenseMonth
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No rows detected!")
		}
		return nil, errors.New("Failed to get the total expense month!" + err.Error())
	}

	return &data, nil

}

func (r *repoTransaction) GetTotalIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeDay, error) {

	query := `
		SELECT DATE_TRUNC('day', date) as day,
		SUM(amount) as total_income
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		GROUP BY day
		ORDER BY day ASC;
	`

	var data model.TotalIncomeDay
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the total income day!" + err.Error())
	}

	return &data, nil

}

func (r *repoTransaction) GetTotalIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeWeek, error) {

	query := `
		SELECT DATE_TRUNC('week', date) as week,
		SUM(amount) as total_income
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		GROUP BY week
		ORDER BY week ASC;
	`

	var data model.TotalIncomeWeek
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the total income week!" + err.Error())
	}

	return &data, nil

}

func (r *repoTransaction) GetTotalIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeMonth, error) {

	query := `
		SELECT DATE_TRUNC('month', date) as month,
		SUM(amount) as total_income
		FROM transactions 
		WHERE user_id = $1 AND type = 'income'
		GROUP BY month
		ORDER BY month ASC;
	`

	var data model.TotalIncomeMonth
	if err := r.db.GetContext(ctx, &data, query, user_id); err != nil {
		return nil, errors.New("Failed to get the total income month!" + err.Error())
	}

	return &data, nil

}
