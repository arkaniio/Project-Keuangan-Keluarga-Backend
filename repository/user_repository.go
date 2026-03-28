package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"project-keuangan-keluarga/model"
)

// ExampleRepository defines the contract for the example data-access layer.
type UserRepository interface {
	CreateNewUser(ctx context.Context, user *model.User) error
}

// exampleRepository is the concrete implementation backed by PostgreSQL via sqlx.
type repoUser struct {
	db *sqlx.DB
}

// NewExampleRepository constructs a new ExampleRepository.
func NewExampleRepository(db *sqlx.DB) UserRepository {
	return &repoUser{db: db}
}

// func to create a new user
func (r *repoUser) CreateNewUser(ctx context.Context, user *model.User) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := r.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transactions")
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users(id, name, email, password, role, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7);
	`

	rows, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.New("Failed to exec the context")
	}

	rows_dected, _ := rows.RowsAffected()
	if rows_dected == 0 {
		return errors.New("Failed to detect the rows affected")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction")
	}

	return nil
}
