package controller

import (
	"context"
	"net/http"
	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
	"time"

	"github.com/google/uuid"
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

	userId, err := middleware.GetTokenId(w, r)
	if err != nil {
		return
	}

	transactions_payload, err := utils.ParsingPayloadTransaction(payload, userId)
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

func (c *ControllerHandlerTransaction) UpdateTransactions_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.UpdatePayloadTransaction
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the json!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	userId, err := middleware.GetTokenId(w, r)
	if err != nil {
		return
	}

	if err := c.TransactionService.UpdateTransaction(ctx, userId, payload); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update the transaction!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully updated the transaction!", nil)

}

func (c *ControllerHandlerTransaction) DeleteTransaction_Bp(w http.ResponseWriter, r *http.Request) {

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	userId, err := middleware.GetTokenId(w, r)
	if err != nil {
		return
	}

	if err := c.TransactionService.DeleteTransaction(ctx, userId); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to delete the transaction!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully deleted the transaction!", nil)

}

func (c *ControllerHandlerTransaction) GetTransactionById_Bp(w http.ResponseWriter, r *http.Request) {

	id_params, err := utils.ParamsChiRouter("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params use chi router!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	transaction_data, err := c.TransactionService.GetTransactionById(ctx, id_params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the transaction by id!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to get the transaction by id!", transaction_data)
}

func (c *ControllerHandlerTransaction) GetAllTransaction_Bp(w http.ResponseWriter, r *http.Request) {

	// Parse pagination from query params
	allowedSorts := []string{"created_at", "amount", "date", "type"}
	params := utils.ParsePaginationParams(r, allowedSorts, "created_at")

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	paginatedData, err := c.TransactionService.GetAllTransaction(ctx, params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get transactions!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully retrieved transactions!", paginatedData)
}

func (c *ControllerHandlerTransaction) GetAvgIncomeDay_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_user_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from middleware token!", err.Error())
		return
	}
	if middleware_user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from middleware token!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	avg_income_data, err := c.TransactionService.GetAvgIncomeDay(ctx, middleware_user_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the income avg data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get avg income p day has been successfully!", avg_income_data)

}

func (c *ControllerHandlerTransaction) GetAvgExpenseDay_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_user_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", err.Error())
		return
	}
	if middleware_user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", false)
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	avg_expense_data, err := c.TransactionService.GetAvgExpenseDay(ctx, middleware_user_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the expense avg data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get avg expense p day has been successfully!", avg_expense_data)

}

func (c *ControllerHandlerTransaction) GetAvgIncomeWeek_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if middleware_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	avg_income_week_data, err := c.TransactionService.GetAvgIncomeWeek(ctx, middleware_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the avg income week data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get avg income p week has been successfully!", avg_income_week_data)

}

func (c *ControllerHandlerTransaction) GetAvgExpenseWeek_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if middleware_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	avg_expense_week_data, err := c.TransactionService.GetAvgExpenseWeek(ctx, middleware_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the avg expense week data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get avg expense p week has been successfully!", avg_expense_week_data)

}

func (c *ControllerHandlerTransaction) GetAvgIncomeMonth_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if middleware_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	avg_income_month_data, err := c.TransactionService.GetAvgIncomeMonth(ctx, middleware_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the avg income month data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get avg income p month has been successfully!", avg_income_month_data)

}

func (c *ControllerHandlerTransaction) GetAvgExpenseMonth_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if middleware_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	avg_expense_month_data, err := c.TransactionService.GetAvgExpenseMonth(ctx, middleware_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the avg expense month data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get avg expense p month has been successfully!", avg_expense_month_data)

}
