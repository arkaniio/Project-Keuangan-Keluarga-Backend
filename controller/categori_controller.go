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

	id_params, err := utils.ParamsChiRouter("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id params!", err.Error())
		return
	}

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the type of uuid!", false)
		return
	}

	if err := c.CategoryService.DeleteCategory(ctx, id_params, middleware_token_id); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to delete the category!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully deleted the category!", nil)

}

func (c *ControllerHandlerCategory) GetCategoryById_Bp(w http.ResponseWriter, r *http.Request) {

	id_params, err := utils.ParamsChiRouter("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params use chi router!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	category_data, err := c.CategoryService.GetCategoryById(ctx, id_params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the category by id!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to get the category by id!", category_data)

}

func (c *ControllerHandlerCategory) GetAllCategory_Bp(w http.ResponseWriter, r *http.Request) {

	// Parse pagination from query params
	allowedSorts := []string{"name", "type"}
	params := utils.ParsePaginationParams(r, allowedSorts, "name")

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	paginatedData, err := c.CategoryService.GetAllCategory(ctx, params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get categories!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully retrieved categories!", paginatedData)

}
