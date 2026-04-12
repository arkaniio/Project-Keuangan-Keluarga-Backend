package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type BudgetService interface {
	CreateNewBudget(ctx context.Context, payload *model.Budget) error
	UpdateBudget(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadBudget) error
	GetBudgetById(ctx context.Context, id uuid.UUID) (*model.Budget, error)
	GetBudgetByUserId(ctx context.Context, user_id uuid.UUID) (*model.Budget, error)
	DeleteBudget(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	GetAllBudget(ctx context.Context, params model.PaginationParams) (model.PaginatedResponse, error)
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

func (s *budgetService) DeleteBudget(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {
	return s.budgetRepo.DeleteBudget(ctx, id, user_id)
}

func (s *budgetService) GetAllBudget(ctx context.Context, params model.PaginationParams) (model.PaginatedResponse, error) {

	budgets, total_items, err := s.budgetRepo.GetAllBudget(ctx, params)
	if err != nil {
		return model.PaginatedResponse{}, errors.New("Failed to get the all budgets!")
	}

	meta := utils.BuildPaginationMeta(total_items, params.Page, params.Limit)

	return model.PaginatedResponse{
		Items:      budgets,
		Pagination: meta,
	}, nil
}
