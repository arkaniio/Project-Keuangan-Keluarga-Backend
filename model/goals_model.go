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
	Priority       string    `db:"priority"`
	Start_date     string    `db:"start_date"`
	End_date       string    `db:"end_date"`
	Is_active      bool      `db:"is_active"`
	Created_at     time.Time `db:"created_at"`
	Updated_at     time.Time `db:"updated_at"`
}

type PayloadGoals struct {
	Id             uuid.UUID `json:"id"`
	User_id        uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	Target_amount  float64   `json:"target_amount"`
	Current_amount float64   `json:"current_amount"`
	Priority       string    `json:"priority"`
	Start_date     string    `json:"start_date"`
	End_date       string    `json:"end_date"`
	Is_active      bool      `json:"is_active"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}

type PayloadUpdateGoals struct {
	Name           *string    `json:"name"`
	Target_amount  *float64   `json:"target_amount"`
	Current_amount *float64   `json:"current_amount"`
	Priority       *string    `json:"priority"`
	Start_date     *string    `json:"start_date"`
	End_date       *string    `json:"end_date"`
	Is_active      *bool      `json:"is_active"`
	Updated_at     *time.Time `json:"updated_at"`
}
