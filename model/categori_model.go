package model

import "github.com/google/uuid"

type Category struct {
	Id     uuid.UUID `db:"id"`
	UserId uuid.UUID `db:"user_id"`
	Name   string    `db:"name"`
	Type   string    `db:"type"`
}

type PayloadCategory struct {
	Id     uuid.UUID `json:"id" validate:"required"`
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Name   string    `json:"name" validate:"required"`
	Type   string    `json:"type" validate:"required"`
}
