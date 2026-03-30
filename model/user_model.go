package model

import (
	"time"

	"github.com/google/uuid"
)

// Example represents a sample entity for demonstrating the clean architecture layers.
// Replace or extend this struct when implementing real domain models.
type User struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Role        string    `db:"role"`
	Profile_img string    `db:"profile_img"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Payload struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required"`
	Role        string    `json:"role" validate:"required"`
	Profile_img string    `json:"profile_img"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PayloadUpdate struct {
	Id          uuid.UUID `json:"id"`
	Name        *string   `json:"name" validate:"required"`
	Email       *string   `json:"email" validate:"required,email"`
	Password    *string   `json:"password" validate:"required"`
	Role        *string   `json:"role" validate:"required"`
	Profile_img *string   `json:"profile_img"`
}
