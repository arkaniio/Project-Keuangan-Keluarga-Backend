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

type PayloadFamilieWithUser struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Created_By uuid.UUID `json:"created_by"`
	User       User      `json:"user"`
	Created_at time.Time `json:"created_at"`
}

type PayloadFamilieWithUserData struct {
	Id         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Created_By uuid.UUID `db:"created_by"`
	Username   string    `db:"username"`
	Email      string    `db:"email"`
	Created_at time.Time `db:"created_at"`
}
