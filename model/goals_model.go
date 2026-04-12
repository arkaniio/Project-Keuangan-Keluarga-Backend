package model

import (
	"time"

	"github.com/google/uuid"
)

type Goals struct {
	Id             uuid.UUID `db:"id"`
	User_id        uuid.UUID `db:"user_id"`
	Name           string    `db:"name"`
	Target_amount  float64   `db:"target_amount"`
	Current_amount float64   `db:"current_amount"`
	Start_date     string    `db:"start_date"`
	Target_date    string    `db:"target_date"`
	Status         string    `db:"status"`
	Created_at     time.Time `db:"created_at"`
	Updated_at     time.Time `db:"updated_at"`
}

type PayloadGoals struct {
	Id             uuid.UUID `json:"id"`
	User_id        uuid.UUID `json:"user_id" validate:"required"`
	Name           string    `json:"name" validate:"required"`
	Target_amount  float64   `json:"target_amount" validate:"required"`
	Current_amount float64   `json:"current_amount" validate:"required"`
	Start_date     string    `json:"start_date" validate:"required"`
	Target_date    string    `json:"target_date" validate:"required"`
	Status         string    `json:"status"`
}

type PayloadUpdateGoals struct {
	Name           *string    `json:"name"`
	Target_amount  *float64   `json:"target_amount"`
	Current_amount *float64   `json:"current_amount"`
	Start_date     *string    `json:"start_date"`
	Target_date    *string    `json:"target_date"`
	Status         *string    `json:"status"`
	Updated_at     *time.Time `json:"updated_at"`
}
