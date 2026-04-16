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

type FamilieRepository interface {
	CreateNewFamilie(ctx context.Context, familie *model.Familie) error
	DeleteFamilie(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	UpdateFamilie(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilie) error
	GetAllFamilie(ctx context.Context, params model.PaginationParams) ([]model.PayloadFamilieWithUser, int, error)
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

	if _, err := tx.ExecContext(ctx, query, familie.Id, familie.Name, familie.Created_By, familie.Created_at); err != nil {
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

func (r *repoFamilie) UpdateFamilie(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilie) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settings the transaction!")
	}

	settings, args, err := utils.UpdateToolsFamilie(payload, user_id)
	if err != nil {
		return errors.New("Failed to settings the update tools query!")
	}

	if _, err := tx.ExecContext(ctx, settings, args...); err != nil {
		return errors.New("Failed to execute the query!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}

func (r *repoFamilie) GetAllFamilie(ctx context.Context, params model.PaginationParams) ([]model.PayloadFamilieWithUser, int, error) {

	// ── Build dynamic WHERE clause ─────────────────────────────
	where := ""
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		where = fmt.Sprintf(" WHERE f.name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// ── Count total items ──────────────────────────────────────
	countQuery := "SELECT COUNT(*) FROM families f JOIN users u ON f.created_by = u.id" + where

	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, errors.New("Failed to count families: " + err.Error())
	}

	// ── Fetch paginated data ───────────────────────────────────
	offset := utils.CalculateOffset(params.Page, params.Limit)

	dataQuery := fmt.Sprintf(`
		SELECT f.id, f.name, f.created_by, u.username, u.email, f.created_at
		FROM families f
		JOIN users u ON f.created_by = u.id
		%s
		ORDER BY f.%s %s
		LIMIT $%d OFFSET $%d
	`, where, params.Sort, params.Order, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	var familie_array []model.PayloadFamilieWithUser
	rows, err := r.db.QueryxContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.New("Failed to get the rows from the db: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var familie_user_data model.PayloadFamilieWithUserData
		if err := rows.StructScan(&familie_user_data); err != nil {
			return nil, 0, errors.New("Failed to get and detect the rows from db: " + err.Error())
		}
		familie_data, err := utils.PayloadJoinDataFamilie(familie_user_data)
		if err != nil {
			return nil, 0, errors.New("Failed to parse familie data: " + err.Error())
		}
		familie_array = append(familie_array, familie_data)
	}

	return familie_array, totalItems, nil

}
