package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"

	"github.com/google/uuid"
)

type BudgetService interface {
	CreateNewBudget(ctx context.Context, payload *model.Budget) error
	UpdateBudget(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadBudget) error
	GetBudgetById(ctx context.Context, id uuid.UUID) (*model.Budget, error)
	GetBudgetByUserId(ctx context.Context, user_id uuid.UUID) (*model.Budget, error)
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

func (s *budgetService) UpdateBudget(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadBudget) error {
	return s.budgetRepo.UpdateBudget(ctx, id, payload)
}

func (s *budgetService) GetBudgetById(ctx context.Context, id uuid.UUID) (*model.Budget, error) {
	return s.budgetRepo.GetBudgetById(ctx, id)
}

func (s *budgetService) GetBudgetByUserId(ctx context.Context, user_id uuid.UUID) (*model.Budget, error) {
	return s.budgetRepo.GetBudgetByUserId(ctx, user_id)
}
