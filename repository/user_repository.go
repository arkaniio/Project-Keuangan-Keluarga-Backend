package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/utils"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateDataUser(id uuid.UUID, ctx context.Context, user model.PayloadUpdate) error
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
		SELECT id, name, email, password, role, profile_img, created_at, updated_at 
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

func (r *repoUser) UpdateDataUser(id uuid.UUID, ctx context.Context, payload model.PayloadUpdate) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := r.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transactions")
	}
	defer tx.Rollback()

	var args []interface{}
	argsID := 1
	var settings []string

	if payload.Name != nil {
		settings = append(settings, fmt.Sprintf("name=$%d", argsID))
		args = append(args, *payload.Name)
		argsID++
	}
	if payload.Email != nil {
		settings = append(settings, fmt.Sprintf("email=$%d", argsID))
		args = append(args, *payload.Email)
		argsID++
	}
	if payload.Password != nil {
		hash_password, err := utils.HashPassword(*payload.Password)
		if err != nil {
			return errors.New("Failed to hash the password")
		}
		settings = append(settings, fmt.Sprintf("password=$%d", argsID))
		args = append(args, hash_password)
		argsID++
	}
	if payload.Role != nil {
		settings = append(settings, fmt.Sprintf("role=$%d", argsID))
		args = append(args, *payload.Role)
		argsID++
	}
	if payload.Profile_img != nil {
		settings = append(settings, fmt.Sprintf("profile_img=$%d", argsID))
		args = append(args, *payload.Profile_img)
		argsID++
	}

	settings = append(settings, fmt.Sprintf("updated_at=$%d", argsID))
	args = append(args, time.Now().UTC())
	argsID++

	full_query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(settings, ", "), argsID)
	args = append(args, id)

	rows, err := tx.ExecContext(ctx, full_query, args...)
	if err != nil {
		return errors.New("Failed to exec the query!" + err.Error())
	}

	result, err := rows.RowsAffected()
	if err != nil {
		return errors.New("Failed to get the rows affected!")
	}

	if result == 0 {
		return errors.New("No one data expected!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to update the query and commit the query!")
	}

	return nil

}
