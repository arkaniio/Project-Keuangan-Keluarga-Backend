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

func (c *ControllerHandlerCategory) UpdateCategory_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.UpdatePayloadCategory
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

	if err := c.CategoryService.UpdateCategory(ctx, userId, payload); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update the category!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully updated the category!", nil)

}

func (c *ControllerHandlerCategory) DeleteCategory_Bp(w http.ResponseWriter, r *http.Request) {

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	userId, err := middleware.GetTokenId(w, r)
	if err != nil {
		return
	}

	if err := c.CategoryService.DeleteCategory(ctx, userId); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to delete the category!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully deleted the category!", nil)

}
