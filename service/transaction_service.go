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
	GetAllTransaction(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
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
		budget_data, err := s.repoBudget.GetActiveBudget(ctx, transactions.UserId)
		if err != nil {
			return errors.New("Failed to get the budget data")
		}
		transactions_total_amount, err := s.repoTransaction.GetTotalExpenseByCategory(ctx, transactions.UserId, transactions.CategoryId)
		if err != nil {
			return errors.New("Failed to get the transactions total for expense!")
		}

		total_amount := transactions_total_amount + transactions.Amount
		if total_amount >= budget_data.Limit_amount {
			return errors.New("total amount must be lower than limit amount")
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

func (s *repoTransactionCombine) GetAllTransaction(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error) {
	items, totalItems, err := s.repoTransaction.GetAllTransaction(ctx, params)
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
	return s.repoTransaction.GetAvgIncomeDay(ctx, user_id)
}

func (s *repoTransactionCombine) GetAvgExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDay, error) {
	return s.repoTransaction.GetAvgExpenseDay(ctx, user_id)
}

func (s *repoTransactionCombine) GetAvgIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeWeek, error) {
	return s.repoTransaction.GetAvgIncomeWeek(ctx, user_id)
}

func (s *repoTransactionCombine) GetAvgExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseWeek, error) {
	return s.repoTransaction.GetAvgExpenseWeek(ctx, user_id)
}

func (s *repoTransactionCombine) GetAvgIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeMonth, error) {
	return s.repoTransaction.GetAvgIncomeMonth(ctx, user_id)
}

func (s *repoTransactionCombine) GetAvgExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseMonth, error) {
	return s.repoTransaction.GetAvgExpenseMonth(ctx, user_id)
}

func (s *repoTransactionCombine) GetTransactionDataInExpenseType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {
	return s.repoTransaction.GetTransactionDataInExpenseType(type_transaction, user_id, ctx)
}

func (s *repoTransactionCombine) GetTransactionDataInIncomeType(type_transaction string, user_id uuid.UUID, ctx context.Context) (*model.Transaction, error) {
	return s.repoTransaction.GetTransactionDataInIncomeType(type_transaction, user_id, ctx)
}

func (s *repoTransactionCombine) GetAvgExpenseDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgExpenseDayNameCategory, error) {
	return s.repoTransaction.GetAvgExpenseDayNameCategory(ctx, user_id)
}

func (s *repoTransactionCombine) GetAvgIncomeDayNameCategory(ctx context.Context, user_id uuid.UUID) (*model.AvgIncomeDayNameCategory, error) {
	return s.repoTransaction.GetAvgIncomeDayNameCategory(ctx, user_id)
}

func (s *repoTransactionCombine) GetTotalExpenseDay(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseDay, error) {
	return s.repoTransaction.GetTotalExpenseDay(ctx, user_id)
}

func (s *repoTransactionCombine) GetTotalExpenseWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseWeek, error) {
	return s.repoTransaction.GetTotalExpenseWeek(ctx, user_id)
}

func (s *repoTransactionCombine) GetTotalExpenseMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalExpenseMonth, error) {
	return s.repoTransaction.GetTotalExpenseMonth(ctx, user_id)
}

func (s *repoTransactionCombine) GetTotalIncomeDay(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeDay, error) {
	return s.repoTransaction.GetTotalIncomeDay(ctx, user_id)
}

func (s *repoTransactionCombine) GetTotalIncomeWeek(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeWeek, error) {
	return s.repoTransaction.GetTotalIncomeWeek(ctx, user_id)
}

func (s *repoTransactionCombine) GetTotalIncomeMonth(ctx context.Context, user_id uuid.UUID) (*model.TotalIncomeMonth, error) {
	return s.repoTransaction.GetTotalIncomeMonth(ctx, user_id)
}
