package model

import "github.com/google/uuid"

type Catregory struct {
	Id     uuid.UUID `db:"id"`
	UserId uuid.UUID `db:"user_id"`
	Name   string    `db:"name"`
	Type   string    `db:"type"`
}

type PayloadCategory struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
}
