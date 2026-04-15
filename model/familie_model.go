package model

import (
	"time"

	"github.com/google/uuid"
)

type Familie struct {
	Id         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Created_By uuid.UUID `db:"created_by"`
	Created_at time.Time `db:"created_at"`
}

type PayloadFamilie struct {
	Name       string    `json:"name"`
	Created_By uuid.UUID `json:"created_by"`
}

type UpdateFamilie struct {
	Name       *string    `json:"name"`
	Created_By *uuid.UUID `json:"created_by"`
}
