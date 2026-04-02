package service

import (
	"project-keuangan-keluarga/repository"
)

type TransactionService interface {
}

type repoTransaction struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &repoTransaction{repo: repo}
}
