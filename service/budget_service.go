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
	GetAllBudget(ctx context.Context, params model.PaginationParams, user_id uuid.UUID) (model.PaginatedResponse, error)
}

type budgetService struct {
	budgetRepo       repository.BudgetRepository
	repoUser         repository.UserRepository
	repoFamilyMember repository.FamilyMemberRepository
}

func NewBudgetService(budgetRepo repository.BudgetRepository, repoUser repository.UserRepository, repoFamilyMember repository.FamilyMemberRepository) BudgetService {
	return &budgetService{budgetRepo: budgetRepo, repoUser: repoUser, repoFamilyMember: repoFamilyMember}
}

func (s *budgetService) CreateNewBudget(ctx context.Context, payload *model.Budget) error {
	users_data, err := s.repoUser.GetUserById(ctx, payload.UserId)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != payload.UserId {
		return errors.New("Failed to access this method, id is not same!")
	}

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, payload.UserId)
	if err != nil {
		return errors.New("User is not a member of any family")
	}
	payload.FamilyMemberId = fm.Id

	return s.budgetRepo.CreateNewBudget(ctx, payload)
}

func (s *budgetService) UpdateBudget(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadBudget) error {

	users_data, err := s.repoUser.GetUserById(ctx, id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != id {
		return errors.New("Failed to access this method, id is not same!")
	}

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

func (s *budgetService) GetAllBudget(ctx context.Context, params model.PaginationParams, user_id uuid.UUID) (model.PaginatedResponse, error) {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return model.PaginatedResponse{}, errors.New("Failed to get the users data!")
	}

	if users_data.Id != user_id {
		return model.PaginatedResponse{}, errors.New("Failed to access this method, id is not same!")
	}

	budgets, total_items, err := s.budgetRepo.GetAllBudget(ctx, params, user_id)
	if err != nil {
		return model.PaginatedResponse{}, errors.New("Failed to get the all budgets!")
	}

	meta := utils.BuildPaginationMeta(total_items, params.Page, params.Limit)

	return model.PaginatedResponse{
		Items:      budgets,
		Pagination: meta,
	}, nil
}
