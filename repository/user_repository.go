package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateDataUser(id uuid.UUID, ctx context.Context, user model.UpdatePayloadUser) error
	GetAllUser(ctx context.Context) ([]model.User, error)
}

type repoUser struct {
	db *sqlx.DB
}

func NewExampleRepository(db *sqlx.DB) UserRepository {
	return &repoUser{db: db}
}

// func to create a new user
func (r *repoUser) CreateNewUser(ctx context.Context, user *model.User) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		db_tx.Rollback()
		return errors.New("Failed to adding the transaction!")
	}

	query := `
		INSERT INTO users(id, username, email, password, role, profile_img, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`

	result, err := db_tx.ExecContext(ctx, query, user.Id, user.Username, user.Email, user.Password, user.Role, user.Profile_img, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.New("Failed to exec the context" + err.Error())
	}

	rows_affected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to get the rows affected!")
	}
	if rows_affected == 0 {
		return errors.New("Failed to get the rows affected!")
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil
}

func (r *repoUser) GetUserByEmail(email string) (*model.User, error) {

	query := `
		SELECT id, username, email, password, role, profile_img, created_at, updated_at 
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
		SELECT id, username, email, password, role, profile_img, created_at, updated_at 
		FROM users WHERE id = $1;
	`

	var user model.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get the data user with id")
		}
		return nil, fmt.Errorf("Failed to get the data user with id")
	}

	return &user, nil

}

func (r *repoUser) UpdateDataUser(id uuid.UUID, ctx context.Context, payload model.UpdatePayloadUser) error {

	db_tx, err := utils.AddTransaction(r.db, ctx)
	if err != nil {
		db_tx.Rollback()
		return errors.New("Failed to get and settings the transactions!")
	}

	full_query, args, err := utils.UpdateToolsUser(payload, id)
	if err != nil {
		return errors.New("Failed to settins the update tools for a user!")
	}

	if _, err := db_tx.ExecContext(ctx, full_query, args...); err != nil {
		return errors.New("Failed to exec the query!" + err.Error())
	}

	if err := db_tx.Commit(); err != nil {
		return errors.New("Failed to update the query and commit the query!")
	}

	return nil

}

func (r *repoUser) GetAllUser(ctx context.Context) ([]model.User, error) {

	query := `
		SELECT id, username, email, role, COALESCE(profile_img, '') as profile_img, created_at, updated_at 
		FROM users;
	`

	var users []model.User
	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, errors.New("Failed to get the all of user data!" + err.Error())
	}

	return users, nil

}
