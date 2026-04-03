package controller

import (
	"context"
	"net/http"
	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
	"time"
)

type ControllerHandlerCategory struct {
	CategoryService service.CategoryService
}

func NewControllerHandlerCategory(categoryService service.CategoryService) *ControllerHandlerCategory {
	return &ControllerHandlerCategory{CategoryService: categoryService}
}

func (c *ControllerHandlerCategory) CreateNewCategory_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.PayloadCategory
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings and decode the json!", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payload", err.Error())
		return
	}

	userId, err := middleware.GetTokenId(w, r)
	if err != nil {
		return
	}

	category, err := utils.ParsingPayloadCategory(payload, userId)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the payload", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.CategoryService.CreateNewCategory(ctx, category); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payload", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to create the new category", nil)

}
