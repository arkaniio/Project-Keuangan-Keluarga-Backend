package model

import (
	"time"

	"github.com/google/uuid"
)

type FamilyMember struct {
	Id        uuid.UUID `db:"id"`
	FamilyId  uuid.UUID `db:"family_id"`
	UserId    uuid.UUID `db:"user_id"`
	Role      string    `db:"role"`
	JoinedAt  time.Time `db:"joined_at"`
}

type PayloadFamilyMember struct {
	FamilyId uuid.UUID `json:"family_id" validate:"required"`
	Role     string    `json:"role" validate:"required"`
}

type UpdateFamilyMember struct {
	FamilyId *uuid.UUID `json:"family_id"`
	Role     *string    `json:"role"`
}

type PayloadFamilyMemberWithUser struct {
	Id       uuid.UUID `json:"id"`
	FamilyId uuid.UUID `json:"family_id"`
	UserId   uuid.UUID `json:"user_id"`
	User     User      `json:"user"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type PayloadFamilyMemberWithUserData struct {
	Id       uuid.UUID `db:"id"`
	FamilyId uuid.UUID `db:"family_id"`
	UserId   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	Role     string    `db:"role"`
	JoinedAt time.Time `db:"joined_at"`
}
