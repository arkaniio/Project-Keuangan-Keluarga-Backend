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

type CategoryRepository interface {
	CreateNewCategory(ctx context.Context, categories *model.Category) error
	UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	GetAllCategory(ctx context.Context, params model.PaginationParams) ([]model.PayloadCategoryWithUser, int, error)
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
		db_tx.Rollback()
		return errors.New("Failed to settings and add the transactions for this method!")
	}

	query := `
		INSERT INTO categories(id, user_id, family_member_id, name, type) 
		VALUES($1, $2, $3, $4, $5);
	`

	rows, err := db_tx.ExecContext(ctx, query, categories.Id, categories.UserId, categories.FamilyMemberId, categories.Name, categories.Type)
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
		db_tx.Rollback()
		return errors.New("Failed to add and settings the transactions")
	}

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

func (r *repoCategory) DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to add and settings the transactions")
	}
	defer db_tx.Rollback()

	query := `
		DELETE FROM categories WHERE id = $1 AND user_id = $2;
	`

	if _, err := db_tx.ExecContext(ctx, query, id, user_id); err != nil {
		return errors.New("Failed to execute the db!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the db!")
	}

	return nil

}

func (r *repoCategory) GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {

	query := `
		SELECT id, user_id, name, type FROM categories WHERE id = $1;
	`

	var category model.Category
	if err := r.db.GetContext(ctx, &category, query, id); err != nil {
		return nil, errors.New("Failed to execute the db!" + err.Error())
	}

	return &category, nil

}

func (r *repoCategory) GetAllCategory(ctx context.Context, params model.PaginationParams) ([]model.PayloadCategoryWithUser, int, error) {

	// ── Build dynamic WHERE clause ─────────────────────────────
	where := ""
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		where = fmt.Sprintf(" WHERE c.name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// ── Count total items ──────────────────────────────────────
	countQuery := "SELECT COUNT(*) FROM categories c JOIN users u ON c.user_id = u.id" + where

	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, errors.New("Failed to count categories: " + err.Error())
	}

	// ── Fetch paginated data ───────────────────────────────────
	offset := utils.CalculateOffset(params.Page, params.Limit)

	dataQuery := fmt.Sprintf(`
		SELECT c.id, c.user_id, c.family_member_id, u.username, u.email, c.name, c.type
		FROM categories c
		JOIN users u ON c.user_id = u.id
		%s
		ORDER BY c.%s %s
		LIMIT $%d OFFSET $%d
	`, where, params.Sort, params.Order, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	var category_array []model.PayloadCategoryWithUser
	rows, err := r.db.QueryxContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.New("Failed to get the rows from the db: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var category_user_data model.PayloadCategoryWithUserData
		if err := rows.StructScan(&category_user_data); err != nil {
			return nil, 0, errors.New("Failed to get and detect the rows from db: " + err.Error())
		}
		category_data, err := utils.PayloadJoinDataCategory(category_user_data)
		if err != nil {
			return nil, 0, errors.New("Failed to parse category data: " + err.Error())
		}
		category_array = append(category_array, category_data)
	}

	return category_array, totalItems, nil

}
