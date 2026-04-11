package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type BudgetService interface {
	CreateNewBudget(ctx context.Context, payload *model.Budget) error
}

type budgetService struct {
	budgetRepo repository.BudgetRepository
}

func NewBudgetService(budgetRepo repository.BudgetRepository) BudgetService {
	return &budgetService{budgetRepo: budgetRepo}
}

func (s *budgetService) CreateNewBudget(ctx context.Context, payload *model.Budget) error {
	return s.budgetRepo.CreateNewBudget(ctx, payload)
}
