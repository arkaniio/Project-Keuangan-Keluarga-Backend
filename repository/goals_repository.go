package repository

import (
	"context"
	"errors"
	"fmt"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GoalsRepository interface {
	CreateNewGoals(ctx context.Context, goals *model.Goals) error
	GetAllGoals(ctx context.Context, params model.PaginationParams, family_id uuid.UUID) ([]model.PayloadGoalsWithUser, int, error)
	DeleteGoals(ctx context.Context, user_id uuid.UUID) error
	UpdateGoals(ctx context.Context, user_id uuid.UUID, payload model.PayloadUpdateGoals) error
	TrackingProgressGoals(ctx context.Context, family_id uuid.UUID) ([]model.ProgressGoals, error)
	RemainingDaysGoals(ctx context.Context, family_id uuid.UUID) ([]model.RemainingDays, error)
}

type repoGoals struct {
	db *sqlx.DB
}

func NewGoalsRepository(db *sqlx.DB) GoalsRepository {
	return &repoGoals{db: db}
}

func (r *repoGoals) CreateNewGoals(ctx context.Context, goals *model.Goals) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settigs the new transaction for goals!")
	}

	query := `
		INSERT INTO goals (id, user_id, family_member_id, name, target_amount, current_amount, start_date, target_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	if _, err := tx.ExecContext(ctx, query, goals.Id, goals.User_id, goals.FamilyMemberId, goals.Name, goals.Target_amount, goals.Current_amount, goals.Start_date, goals.Target_date, goals.Status, goals.Created_at, goals.Updated_at); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query for goals!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction for goals!")
	}

	return nil

}

func (r *repoGoals) GetAllGoals(ctx context.Context, params model.PaginationParams, family_id uuid.UUID) ([]model.PayloadGoalsWithUser, int, error) {

	args := []interface{}{family_id}
	argIdx := 2

	where := " WHERE fm.family_id = $1"

	if params.Search != "" {
		where += fmt.Sprintf(" AND g.name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	total_count_query := `
		SELECT COUNT(*) 
		FROM goals g 
		JOIN family_members fm ON g.family_member_id = fm.id
		JOIN users u ON g.user_id = u.id` + where

	var total_items int
	if err := r.db.GetContext(ctx, &total_items, total_count_query, args...); err != nil {
		return nil, 0, errors.New("Failed to get the total items for get total count: " + err.Error())
	}

	offset := utils.CalculateOffset(params.Page, params.Limit)

	query := fmt.Sprintf(`
		SELECT g.id, g.user_id, g.family_member_id, u.username, u.email, COALESCE(u.profile_img, '') as profile_img, g.name, g.target_amount, g.current_amount, g.start_date, g.target_date, g.status, g.created_at, g.updated_at
		FROM goals g
		JOIN family_members fm ON g.family_member_id = fm.id
		LEFT JOIN users u ON g.user_id = u.id
		%s
		ORDER BY g.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	var goals_data []model.PayloadGoalsWithUser
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.New("Failed to get the goals data: " + err.Error())
	}

	for rows.Next() {
		var goals model.PayloadGoalsWithUserData
		if err := rows.StructScan(&goals); err != nil {
			return nil, 0, errors.New("Failed to scan the goals data")
		}
		goals_struct, err := utils.PayloadJoinDataGoals(goals)
		if err != nil {
			return nil, 0, errors.New("Failed to parse the goals data")
		}
		goals_data = append(goals_data, *goals_struct)
	}

	return goals_data, total_items, nil

}

func (r *repoGoals) DeleteGoals(ctx context.Context, id uuid.UUID) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settigs the new transaction for goals!")
	}

	query := `
		DELETE FROM goals WHERE user_id = $1
	`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query for goals!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction for goals!")
	}

	return nil

}

func (r *repoGoals) UpdateGoals(ctx context.Context, user_id uuid.UUID, payload model.PayloadUpdateGoals) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settigs the new transaction for goals!")
	}

	query, args, err := utils.UpdateToolsGoals(payload, user_id)
	if err != nil {
		return errors.New("Failed to update the goals!")
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query for goals!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction for goals!")
	}

	return nil

}

func (r *repoGoals) TrackingProgressGoals(ctx context.Context, family_id uuid.UUID) ([]model.ProgressGoals, error) {

	query := `
		SELECT 
    	g.id,
    	g.name,
    	g.target_amount,
    	g.current_amount,
    	g.target_date
		FROM goals g
		JOIN family_members fm ON g.family_member_id = fm.id
		WHERE fm.family_id = $1;
	`

	rows, err := r.db.QueryxContext(ctx, query, family_id)
	if err != nil {
		return nil, errors.New("Failed to get the rows result from db: " + err.Error())
	}

	var goals_data []model.ProgressGoals

	for rows.Next() {
		var g model.ProgressGoals
		if err := rows.StructScan(&g); err != nil {
			return nil, errors.New("Failed to scan the struct progress goals!")
		}

		g.Progress = g.Current_amount / g.Target_amount * 100
		goals_data = append(goals_data, g)
	}

	return goals_data, nil

}

func (r *repoGoals) RemainingDaysGoals(ctx context.Context, family_id uuid.UUID) ([]model.RemainingDays, error) {

	query := `
		SELECT 
    	g.id,
    	g.name,
    	g.target_amount,
    	g.current_amount,
    	g.target_date
		FROM goals g
		JOIN family_members fm ON g.family_member_id = fm.id
		WHERE fm.family_id = $1;
	`

	rows, err := r.db.QueryxContext(ctx, query, family_id)
	if err != nil {
		return nil, errors.New("Failed to get the rows of the result in db: " + err.Error())
	}

	var goals_data []model.RemainingDays

	for rows.Next() {
		var g model.RemainingDays
		if err := rows.StructScan(&g); err != nil {
			return nil, errors.New("Failed to scan the struct in model!")
		}

		g.Remaining_Days = int(time.Until(g.Target_Date).Hours() / 24)

		goals_data = append(goals_data, g)
	}

	return goals_data, nil

}
