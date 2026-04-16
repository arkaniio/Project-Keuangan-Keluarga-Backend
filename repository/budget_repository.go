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

type BudgetRepository interface {
	CreateNewBudget(ctx context.Context, payload *model.Budget) error
	UpdateBudget(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadBudget) error
	GetBudgetById(ctx context.Context, id uuid.UUID) (*model.Budget, error)
	GetBudgetByUserId(ctx context.Context, family_id uuid.UUID) (*model.Budget, error)
	GetActiveBudget(ctx context.Context, family_id uuid.UUID) (*model.Budget, error)
	DeleteBudget(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	GetAllBudget(ctx context.Context, params model.PaginationParams, family_id uuid.UUID) ([]model.BudgetWithCategoryAndUserData, int, error)
}

type repoBudget struct {
	db *sqlx.DB
}

func NewBudgetRepository(db *sqlx.DB) BudgetRepository {
	return &repoBudget{db: db}
}

func (r *repoBudget) CreateNewBudget(ctx context.Context, payload *model.Budget) error {

	query := `
		INSERT INTO budgets (id, user_id, family_member_id, category_id, limit_amount, period, start_date, end_date, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query, payload.Id, payload.UserId, payload.FamilyMemberId, payload.Category_Id, payload.Limit_amount, payload.Period, payload.StartDate, payload.EndDate, payload.IsActive, payload.Created_at)
	if err != nil {
		return err
	}

	return nil

}

func (r *repoBudget) UpdateBudget(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadBudget) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to get the transactions!")
	}

	updateQuery, args, err := utils.UpdateToolsBudget(payload, id)
	if err != nil {
		return errors.New("Failed to update the query!")
	}

	if _, err := tx.ExecContext(ctx, updateQuery, args...); err != nil {
		return errors.New("Failed to execute the query!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}

func (r *repoBudget) GetBudgetById(ctx context.Context, id uuid.UUID) (*model.Budget, error) {

	query := `
		SELECT id, user_id, category_id, limit_amount, period, start_date, end_date, is_active
		FROM budgets
		WHERE id = $1
	`

	var budget model.Budget
	if err := r.db.GetContext(ctx, &budget, query, id); err != nil {
		return nil, err
	}

	return &budget, nil

}

func (r *repoBudget) GetBudgetByUserId(ctx context.Context, family_id uuid.UUID) (*model.Budget, error) {

	query := `
		SELECT b.id, b.user_id, b.category_id, b.limit_amount, b.period, b.start_date, b.end_date, b.is_active
		FROM budgets b
		JOIN family_members fm ON b.family_member_id = fm.id
		WHERE fm.family_id = $1
	`

	var budget model.Budget
	if err := r.db.GetContext(ctx, &budget, query, family_id); err != nil {
		return nil, err
	}

	return &budget, nil

}

func (r *repoBudget) GetActiveBudget(ctx context.Context, family_id uuid.UUID) (*model.Budget, error) {

	query := `
		SELECT b.id, b.user_id, b.category_id, b.limit_amount, b.period, b.start_date, b.end_date, b.is_active
		FROM budgets b
		JOIN family_members fm ON b.family_member_id = fm.id
		WHERE fm.family_id = $1 AND b.is_active = true
	`

	var budget model.Budget
	if err := r.db.GetContext(ctx, &budget, query, family_id); err != nil {
		return nil, err
	}

	return &budget, nil

}

func (r *repoBudget) DeleteBudget(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to get and add the transaction!")
	}

	query := `
		DELETE FROM budgets WHERE id = $1 AND user_id = $2;
	`

	if _, err := tx.ExecContext(ctx, query, id, user_id); err != nil {
		return errors.New("Failed to execute the context and failed to update using query!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}

func (r *repoBudget) GetAllBudget(ctx context.Context, params model.PaginationParams, family_id uuid.UUID) ([]model.BudgetWithCategoryAndUserData, int, error) {

	args := []interface{}{family_id}
	argIdx := 2

	where := " WHERE fm.family_id = $1"

	if params.Search != "" {
		where += fmt.Sprintf(" AND b.period ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	countQuery := "SELECT COUNT(*) FROM budgets b JOIN family_members fm ON b.family_member_id = fm.id JOIN categories c ON b.category_id = c.id JOIN users u ON b.user_id = u.id" + where

	var total_items int
	if err := r.db.GetContext(ctx, &total_items, countQuery, args...); err != nil {
		return nil, 0, errors.New("Failed to counting the data!")
	}

	offset := utils.CalculateOffset(params.Page, params.Limit)

	query := fmt.Sprintf(`
		SELECT b.id, b.user_id, b.family_member_id, u.username as user_name, u.email as user_email, c.name as category_name, c.type as category_type, b.limit_amount, b.period, b.start_date, b.end_date, b.is_active, b.created_at, b.updated_at
		FROM budgets b
		JOIN family_members fm ON b.family_member_id = fm.id
		JOIN categories c ON b.category_id = c.id
		JOIN users u ON b.user_id = u.id
		%s
		ORDER BY b.%s %s
		LIMIT $%d OFFSET $%d
	`, where, params.Sort, params.Order, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	var budget_array []model.BudgetWithCategoryAndUserData
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.New("Failed to counting the data!")
	}

	for rows.Next() {
		var budget_data model.BudgetWithCategoryAndUser
		if err := rows.StructScan(&budget_data); err != nil {
			return nil, 0, errors.New("Failed to struct the budget type model")
		}
		budgets, err := utils.PayloadJoinDataCategoryAndUser(budget_data)
		if err != nil {
			return nil, 0, errors.New("Failed to get the budgets!")
		}
		budget_array = append(budget_array, *budgets)
	}

	return budget_array, total_items, nil

}
