package repository

import (
	"context"
	"project-keuangan-keluarga/model"

	"github.com/jmoiron/sqlx"
)

type BudgetRepository interface {
	CreateNewBudget(ctx context.Context, payload *model.Budget) error
}

type repoBudget struct {
	db *sqlx.DB
}

func NewBudgetRepository(db *sqlx.DB) BudgetRepository {
	return &repoBudget{db: db}
}

func (r *repoBudget) CreateNewBudget(ctx context.Context, payload *model.Budget) error {

	query := `
		INSERT INTO budgets (id, user_id, category_id, limit_amount, period, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query, payload.Id, payload.UserId, payload.Category_Id, payload.Limit_amount, payload.Period, payload.StartDate, payload.EndDate, payload.IsActive)
	if err != nil {
		return err
	}

	return nil

}
