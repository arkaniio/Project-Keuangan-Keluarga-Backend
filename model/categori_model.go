package model

import "github.com/google/uuid"

type Category struct {
	Id     uuid.UUID `db:"id"`
	UserId uuid.UUID `db:"user_id"`
	Name   string    `db:"name"`
	Type   string    `db:"type"`
}

type PayloadCategory struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	Name   string    `json:"name" validate:"required"`
	Type   string    `json:"type" validate:"required"`
}

type UpdatePayloadCategory struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
}

type PayloadCategoryWithUser struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	User   User      `json:"user"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
}

type PayloadCategoryWithUserData struct {
	Id       uuid.UUID `json:"id"`
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
}
