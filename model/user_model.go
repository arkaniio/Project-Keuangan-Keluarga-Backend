package model

import (
	"time"

	"github.com/google/uuid"
)

// Example represents a sample entity for demonstrating the clean architecture layers.
// Replace or extend this struct when implementing real domain models.
type User struct {
	Id          uuid.UUID `db:"id"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Role        string    `db:"role"`
	Profile_img string    `db:"profile_img"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Payload struct {
	Id          uuid.UUID `json:"id"`
	Username    string    `json:"username" validate:"required"`
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

type UpdatePayloadUser struct {
	Id          uuid.UUID `json:"id"`
	Username    *string   `json:"username"`
	Email       *string   `json:"email"`
	Password    *string   `json:"password"`
	Role        *string   `json:"role"`
	Profile_img *string   `json:"profile_img"`
}
