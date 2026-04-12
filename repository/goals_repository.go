package repository

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"

	"github.com/jmoiron/sqlx"
)

type GoalsRepository interface {
	CreateNewGoals(ctx context.Context, goals *model.Goals) error
}

type repoGoals struct {
	db *sqlx.DB
}

func NewGoalsRepository(db *sqlx.DB) repoGoals {
	return repoGoals{db: db}
}

func (r *repoGoals) CreateNewGoals(ctx context.Context, goals *model.Goals) error {

	tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		return errors.New("Failed to settigs the new transaction for goals!")
	}

	query := `
		INSERT INTO goals (id, user_id, name, target_amount, current_amount, start_date, target_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	if _, err := tx.ExecContext(ctx, query, goals.Id, goals.User_id, goals.Name, goals.Target_amount, goals.Current_amount, goals.Start_date, goals.Target_date, goals.Status, goals.Created_at, goals.Updated_at); err != nil {
		tx.Rollback()
		return errors.New("Failed to execute the query for goals!")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.New("Failed to commit the transaction for goals!")
	}

	return nil

}
