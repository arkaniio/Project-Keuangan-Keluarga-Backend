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

type BudgetController struct {
	budgetService service.BudgetService
}

func NewBudgetController(budgetService service.BudgetService) *BudgetController {
	return &BudgetController{budgetService: budgetService}
}

func (c *BudgetController) CreateNewBudget_Bp(w http.ResponseWriter, r *http.Request) {

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
