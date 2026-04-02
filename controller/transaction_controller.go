package controller

import (
	"project-keuangan-keluarga/service"
)

type ControllerHandlerTransaction struct {
	TransactionService service.TransactionService
}

func NewControllerHandlerTransaction(transactionService service.TransactionService) *ControllerHandlerTransaction {
	return &ControllerHandlerTransaction{TransactionService: transactionService}
}
