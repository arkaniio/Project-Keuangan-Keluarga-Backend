package model

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	Id           uuid.UUID `db:"id"`
	UserId       uuid.UUID `db:"user_id"`
	Category_Id  uuid.UUID `db:"category_id"`
	Limit_amount int64     `db:"limit_amount"`
	Period       string    `db:"period"`
	StartDate    time.Time `db:"start_date"`
	EndDate      time.Time `db:"end_date"`
	IsActive     bool      `db:"is_active"`
	Created_at   time.Time `db:"created_at"`
	Updated_at   time.Time `db:"updated_at"`
}

type PayloadBudget struct {
	Category_Id  uuid.UUID `json:"category_id" validate:"required"`
	Limit_amount int64     `json:"limit_amount" validate:"required"`
	Period       string    `json:"period" validate:"required"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_date" validate:"required"`
	IsActive     bool      `json:"is_active" validate:"required"`
}

type UpdatePayloadBudget struct {
	Id           uuid.UUID  `json:"id"`
	Category_Id  *uuid.UUID `json:"category_id"`
	Limit_amount *int64     `json:"limit_amount"`
	Period       *string    `json:"period"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsActive     *bool      `json:"is_active"`
}

type GetLimitAmountByNameCategory struct {
	Category_Name string `db:"category_name"`
	Limit_amount  int64  `db:"limit_amount"`
}

type BudgetWithCategoryAndUser struct {
	Id            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"user_id"`
	Category_Id   uuid.UUID `json:"category_id"`
	Limit_amount  int64     `json:"limit_amount"`
	Period        string    `json:"period"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsActive      bool      `json:"is_active"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	Category_Name string    `json:"category_name"`
	Category_Type string    `json:"category_type"`
	User_Name     string    `json:"user_name"`
	User_Email    string    `json:"user_email"`
}

type BudgetWithCategoryAndUserData struct {
	Id            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"user_id"`
	User_Data     User      `json:"user"`
	Category_Id   uuid.UUID `json:"category_id"`
	Category_Data Category  `json:"category"`
	Limit_amount  int64     `json:"limit_amount"`
	Period        string    `json:"period"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsActive      bool      `json:"is_active"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}
