package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type TransactionService interface {
	CreateNewTransactions(ctx context.Context, transactions *model.Transaction) error
	UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetTransactionById(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetAllTransaction(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
	GetAvgIncomeDay(ctx context.Context, user_id uuid.UUID) ([]model.AvgIncomeDay, error)
	GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) ([]model.AvgExpenseDay, error)
}

type repoTransaction struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &repoTransaction{repo: repo}
}

func (s *repoTransaction) CreateNewTransactions(ctx context.Context, transactions *model.Transaction) error {
	return s.repo.CreateNewTransactions(ctx, transactions)
}

func (s *repoTransaction) UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error {
	return s.repo.UpdateTransaction(ctx, id, payload)
}

func (s *repoTransaction) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTransaction(ctx, id)
}

func (s *repoTransaction) GetTransactionById(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	return s.repo.GetTransactionById(ctx, id)
}

func (s *repoTransaction) GetAllTransaction(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error) {
	items, totalItems, err := s.repo.GetAllTransaction(ctx, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}

func (s *repoTransaction) GetAvgIncomeDay(ctx context.Context, user_id uuid.UUID) ([]model.AvgIncomeDay, error) {
	return s.repo.GetAvgIncomeDay(ctx, user_id)
}

func (s *repoTransaction) GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) ([]model.AvgExpenseDay, error) {
	return s.repo.GetAvgExpenseDay(ctx, user_id)
}
