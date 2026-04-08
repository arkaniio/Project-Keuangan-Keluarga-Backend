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
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Type        string    `json:"type" validate:"required"`
	Amount      int64     `json:"amount" validate:"required"`
	CategoryId  uuid.UUID `json:"category_id" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

type UpdatePayloadTransaction struct {
	Type        *string    `json:"type"`
	Amount      *int64     `json:"amount"`
	CategoryId  *uuid.UUID `json:"category_id"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}

type PayloadTransactionWithCategory struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	CategoryId  uuid.UUID `json:"category_id"`
	Category    Category  `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PayloadTransactionDataCategory struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	CategoryId  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AvgIncomeDay struct {
	Day       string  `db:"day"`
	AvgIncome float64 `db:"avg_income"`
}

type AvgExpenseDay struct {
	Day        string  `db:"day"`
	AvgExpense float64 `db:"avg_expense"`
}

type AvgIncomeWeek struct {
	Week      string  `db:"week"`
	AvgIncome float64 `db:"income"`
}

type AvgExpenseWeek struct {
	Week       string  `db:"week"`
	AvgExpense float64 `db:"expense"`
}

type AvgIncomeMonth struct {
	Month     string  `db:"month"`
	AvgIncome float64 `db:"income"`
}

type AvgExpenseMonth struct {
	Month      string  `db:"month"`
	AvgExpense float64 `db:"expense"`
}

type AvgExpenseDayNameCategory struct {
	Day        string  `db:"day"`
	Name       string  `db:"name"`
	AvgExpense float64 `db:"avg_expense"`
}

type AvgIncomeDayNameCategory struct {
	Day       string  `db:"day"`
	Category  string  `db:"category"`
	AvgIncome float64 `db:"avg_income"`
}
