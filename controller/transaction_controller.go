package controller

import (
	"context"
	"net/http"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
	"time"
)

type ControllerHandlerTransaction struct {
	TransactionService service.TransactionService
}

func NewControllerHandlerTransaction(transactionService service.TransactionService) *ControllerHandlerTransaction {
	return &ControllerHandlerTransaction{TransactionService: transactionService}
}

func (c *ControllerHandlerTransaction) CreateNewTransactions_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.PayloadTransaction
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the payload and decode!", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payloads!", err.Error())
		return
	}

	transactions_payload, err := utils.ParsingPayloadTransaction(payload)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parse the payload!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.TransactionService.CreateNewTransactions(ctx, transactions_payload); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create the transactions!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully created the transactions!", nil)

}
