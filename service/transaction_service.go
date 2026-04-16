package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type TransactionService interface {
	CreateNewTransactions(ctx context.Context, transactions *model.Transaction) error
	UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error
	GetTransactionByUserId(ctx context.Context, user_id uuid.UUID) (*model.Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetTransactionById(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetAllTransaction(ctx context.Context, user_id uuid.UUID, params model.PaginationParams) (*model.PaginatedResponse, error)
	GetAvgIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDay, error)
	GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDay, error)
	GetAvgIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeWeek, error)
	GetAvgExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseWeek, error)
	GetAvgIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeMonth, error)
	GetAvgExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseMonth, error)
	GetTransactionDataInExpenseType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error)
	GetTransactionDataInIncomeType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error)
	GetAvgExpenseDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDayNameCategory, error)
	GetAvgIncomeDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDayNameCategory, error)
	GetTotalExpenseByCategory(ctx context.Context, user_id uuid.UUID, category_id uuid.UUID) (int64, error)
	GetTotalExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseDay, error)
	GetTotalExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseWeek, error)
	GetTotalExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseMonth, error)
	GetTotalIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeDay, error)
	GetTotalIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeWeek, error)
	GetTotalIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeMonth, error)
}

type repoTransactionCombine struct {
	repoTransaction  repository.TransactionRepository
	repoBudget       repository.BudgetRepository
	repoFamilyMember repository.FamilyMemberRepository
}

func NewTransactionService(repoTransaction repository.TransactionRepository, repoBudget repository.BudgetRepository, repoFamilyMember repository.FamilyMemberRepository) TransactionService {
	return &repoTransactionCombine{repoTransaction: repoTransaction, repoBudget: repoBudget, repoFamilyMember: repoFamilyMember}
}

func (s *repoTransactionCombine) CreateNewTransactions(ctx context.Context, transactions *model.Transaction) error {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, transactions.UserId)
	if err != nil {
		return errors.New("User is not a member of any family")
	}
	transactions.FamilyMemberId = fm.Id

	if transactions.Type != "income" && transactions.Type != "expense" {
		return errors.New("Failed to detect for a type, invalid type!")
	}

	if transactions.Type == "expense" {
		budget_data, err := s.repoBudget.GetActiveBudget(ctx, fm.FamilyId)
		if err == nil && budget_data != nil {
			transactions_total_amount, err := s.repoTransaction.GetTotalExpenseByCategory(ctx, fm.FamilyId, transactions.CategoryId)
			if err != nil {
				return errors.New("Failed to get the transactions total for expense!")
			}

			total_amount := transactions_total_amount + transactions.Amount
			if total_amount >= budget_data.Limit_amount {
				return errors.New("total amount must be lower than limit amount")
			}
		}
	}

	if err := s.repoTransaction.CreateNewTransactions(ctx, transactions); err != nil {
		return errors.New("Failed to create new transactions!")
	}

	return nil

}

func (s *repoTransactionCombine) GetTotalExpenseByCategory(ctx context.Context, user_id uuid.UUID, category_id uuid.UUID) (int64, error) {
	return s.repoTransaction.GetTotalExpenseByCategory(ctx, user_id, category_id)
}

func (s *repoTransactionCombine) UpdateTransaction(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadTransaction) error {
	return s.repoTransaction.UpdateTransaction(ctx, id, payload)
}

func (s *repoTransactionCombine) GetTransactionByUserId(ctx context.Context, user_id uuid.UUID) (*model.Transaction, error) {
	return s.repoTransaction.GetTransactionByUserId(ctx, user_id)
}

func (s *repoTransactionCombine) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.repoTransaction.DeleteTransaction(ctx, id)
}

func (s *repoTransactionCombine) GetTransactionById(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	return s.repoTransaction.GetTransactionById(ctx, id)
}

func (s *repoTransactionCombine) GetAllTransaction(ctx context.Context, user_id uuid.UUID, params model.PaginationParams) (*model.PaginatedResponse, error) {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return nil, err
	}

	if fm == nil {
		return &model.PaginatedResponse{
			Items: []model.PayloadTransactionWithCategory{},
			Pagination: model.PaginationMeta{
				TotalItems:   0,
				TotalPages:   0,
				CurrentPage:  params.Page,
				PerPage:      params.Limit,
			},
		}, nil
	}

	items, totalItems, err := s.repoTransaction.GetAllTransaction(ctx, fm.FamilyId, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}

func (s *repoTransactionCombine) GetAvgIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDay, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgIncomeDay{}, nil
	}
	return s.repoTransaction.GetAvgIncomeDay(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDay, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgExpenseDay{}, nil
	}
	return s.repoTransaction.GetAvgExpenseDay(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetAvgIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeWeek, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgIncomeWeek{}, nil
	}
	return s.repoTransaction.GetAvgIncomeWeek(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetAvgExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseWeek, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgExpenseWeek{}, nil
	}
	return s.repoTransaction.GetAvgExpenseWeek(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetAvgIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeMonth, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgIncomeMonth{}, nil
	}
	return s.repoTransaction.GetAvgIncomeMonth(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetAvgExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseMonth, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgExpenseMonth{}, nil
	}
	return s.repoTransaction.GetAvgExpenseMonth(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTransactionDataInExpenseType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return nil, nil
	}
	return s.repoTransaction.GetTransactionDataInExpenseType(type_transaction, fm.FamilyId, ctx)
}

func (s *repoTransactionCombine) GetTransactionDataInIncomeType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return nil, nil
	}
	return s.repoTransaction.GetTransactionDataInIncomeType(type_transaction, fm.FamilyId, ctx)
}

func (s *repoTransactionCombine) GetAvgExpenseDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDayNameCategory, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgExpenseDayNameCategory{}, nil
	}
	return s.repoTransaction.GetAvgExpenseDayNameCategory(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetAvgIncomeDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDayNameCategory, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.AvgIncomeDayNameCategory{}, nil
	}
	return s.repoTransaction.GetAvgIncomeDayNameCategory(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTotalExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseDay, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.TotalExpenseDay{}, nil
	}
	return s.repoTransaction.GetTotalExpenseDay(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTotalExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseWeek, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.TotalExpenseWeek{}, nil
	}
	return s.repoTransaction.GetTotalExpenseWeek(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTotalExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseMonth, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.TotalExpenseMonth{}, nil
	}
	return s.repoTransaction.GetTotalExpenseMonth(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTotalIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeDay, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.TotalIncomeDay{}, nil
	}
	return s.repoTransaction.GetTotalIncomeDay(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTotalIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeWeek, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.TotalIncomeWeek{}, nil
	}
	return s.repoTransaction.GetTotalIncomeWeek(ctx, fm.FamilyId)
}

func (s *repoTransactionCombine) GetTotalIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeMonth, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return &model.TotalIncomeMonth{}, nil
	}
	return s.repoTransaction.GetTotalIncomeMonth(ctx, fm.FamilyId)
}
