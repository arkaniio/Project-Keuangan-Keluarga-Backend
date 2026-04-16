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

type FamilyMemberRepository interface {
	CreateFamilyMember(ctx context.Context, member *model.FamilyMember) error
	UpdateFamilyMember(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilyMember) error
	DeleteFamilyMember(ctx context.Context, user_id uuid.UUID) error
	GetAllFamilyMember(ctx context.Context, params model.PaginationParams) ([]model.PayloadFamilyMemberWithUser, int, error)
}

type repoFamilyMember struct {
	db *sqlx.DB
}

func NewFamilyMemberRepository(db *sqlx.DB) FamilyMemberRepository {
	return &repoFamilyMember{db: db}
}

func (r *repoFamilyMember) CreateFamilyMember(ctx context.Context, member *model.FamilyMember) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settings the transaction!")
	}

	query := `
		INSERT INTO family_members(id, family_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4, $5);
	`

	if _, err := tx.ExecContext(ctx, query, member.Id, member.FamilyId, member.UserId, member.Role, member.JoinedAt); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}

func (r *repoFamilyMember) UpdateFamilyMember(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilyMember) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settings the transaction!")
	}

	settings, args, err := utils.UpdateToolsFamilyMember(payload, user_id)
	if err != nil {
		return errors.New("Failed to settings the update tools query!")
	}

	if _, err := tx.ExecContext(ctx, settings, args...); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}

func (r *repoFamilyMember) DeleteFamilyMember(ctx context.Context, user_id uuid.UUID) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settings the transaction!")
	}

	query := `
		DELETE FROM family_members WHERE user_id = $1;
	`

	if _, err := tx.ExecContext(ctx, query, user_id); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction!")
	}

	return nil

}

func (r *repoFamilyMember) GetAllFamilyMember(ctx context.Context, params model.PaginationParams) ([]model.PayloadFamilyMemberWithUser, int, error) {

	// ── Build dynamic WHERE clause ─────────────────────────────
	where := ""
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		where = fmt.Sprintf(" WHERE fm.role ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// ── Count total items ──────────────────────────────────────
	countQuery := "SELECT COUNT(*) FROM family_members fm JOIN users u ON fm.user_id = u.id" + where

	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, errors.New("Failed to count family members: " + err.Error())
	}

	// ── Fetch paginated data ───────────────────────────────────
	offset := utils.CalculateOffset(params.Page, params.Limit)

	dataQuery := fmt.Sprintf(`
		SELECT fm.id, fm.family_id, fm.user_id, u.username, u.email, fm.role, fm.joined_at
		FROM family_members fm
		JOIN users u ON fm.user_id = u.id
		%s
		ORDER BY fm.%s %s
		LIMIT $%d OFFSET $%d
	`, where, params.Sort, params.Order, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	var member_array []model.PayloadFamilyMemberWithUser
	rows, err := r.db.QueryxContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.New("Failed to get the rows from the db: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var member_user_data model.PayloadFamilyMemberWithUserData
		if err := rows.StructScan(&member_user_data); err != nil {
			return nil, 0, errors.New("Failed to get and detect the rows from db: " + err.Error())
		}
		member_data, err := utils.PayloadJoinDataFamilyMember(member_user_data)
		if err != nil {
			return nil, 0, errors.New("Failed to parse family member data: " + err.Error())
		}
		member_array = append(member_array, member_data)
	}

	return member_array, totalItems, nil

}
