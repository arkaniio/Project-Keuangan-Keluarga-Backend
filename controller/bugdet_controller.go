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

type ControllerBudget struct {
	budgetService service.BudgetService
}

func NewBudgetController(budgetService service.BudgetService) *ControllerBudget {
	return &ControllerBudget{budgetService: budgetService}
}

func (c *ControllerBudget) CreateNewBudget_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id from middleware!", err.Error())
		return
	}
	if middleware_token == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid!", false)
		return
	}

	var payload model.PayloadBudget
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the payload budget!", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payloads!", err.Error())
		return
	}

	budgets, err := utils.ParsingPayloadBudget(payload, middleware_token)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the payload budget!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.budgetService.CreateNewBudget(ctx, budgets); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create budget per user!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusBadRequest, "Success to create new budget data!", true)

}

func (c *ControllerBudget) UpdateBudget_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id from middleware!", err.Error())
		return
	}
	if middleware_token == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid!", false)
		return
	}

	var payload model.UpdatePayloadBudget
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the payload budget!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.budgetService.UpdateBudget(ctx, middleware_token, payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update budget per user!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusBadRequest, "Success to update budget data!", true)

}

func (c *ControllerBudget) DeleteBudget_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware id from token!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", false)
		return
	}

	id_params, err := utils.ParamsChiRouter("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the id params!", err.Error())
		return
	}
	if id_params == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid type for params!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.budgetService.DeleteBudget(ctx, id_params, middleware_token_id); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete tthe budget based on id and user_id!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to delete the budget!", true)

}

func (c *ControllerBudget) GetAllBudget_Bp(w http.ResponseWriter, r *http.Request) {

	allowed_sort := []string{"created_at", "limit_amount"}
	parsing_params := utils.ParsePaginationParams(r, allowed_sort, "created_at")

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	budgets_data, err := c.budgetService.GetAllBudget(ctx, parsing_params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the budgets data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to get all of the data budgets!", budgets_data)

}
