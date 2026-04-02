package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type TransactionService interface {
	CreateNewTransactions(ctx context.Context, transactions *model.Transaction) error
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
