package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"project-keuangan-keluarga/model"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type repoUser struct {
	db *sqlx.DB
}

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
		INSERT INTO users(id, name, email, password, role, profile_img, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`

	rows, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Role, user.Profile_img, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.New("Failed to exec the context" + err.Error())
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

func (r *repoUser) GetUserByEmail(email string) (*model.User, error) {

	query := `
		SELECT id, name, email, password, role, profile_img, created_at, updated_at 
		FROM users WHERE email = $1;
	`

	var user model.User
	if err := r.db.Get(&user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get the email!")
	}

	return &user, nil

}

func (r *repoUser) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {

	query := `
		SELECT id, name, email, password, role, created_at, updated_at 
		FROM users WHERE id = $1;
	`

	if err := r.db.GetContext(ctx, &model.User{}, query, id); err != nil {
		return nil, err
	}

	return nil, nil

}
