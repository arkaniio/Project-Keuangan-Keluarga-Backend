package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id          uuid.UUID `db:"id"`
	UserId      uuid.UUID `db:"user_id"`
	Type        string    `db:"type"`
	Amount      int64     `db:"amount"`
	CategoryId  uuid.UUID `db:"category_id"`
	Description string    `db:"description"`
	Date        time.Time `db:"date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PayloadTransaction struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	UserId      uuid.UUID `json:"user_id" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	Amount      int64     `json:"amount" validate:"required"`
	CategoryId  uuid.UUID `json:"category_id" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}
