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

type ControllerHandlerFamilie struct {
	FamilieService service.FamilieService
}

func NewControllerHandlerFamilie(familieService service.FamilieService) ControllerHandlerFamilie {
	return ControllerHandlerFamilie{FamilieService: familieService}
}

func (c *ControllerHandlerFamilie) CreateNewFamilie_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.PayloadFamilie
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the json!", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payload!", err.Error())
		return
	}

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	familie, err := utils.ParsingPayloadFamilie(payload, middleware_token_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the payload!", err.Error())
		return
	}

	if err := c.FamilieService.CreateNewFamilie(ctx, familie); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create the new familie!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully created the new familie!", nil)

}

func (c *ControllerHandlerFamilie) DeleteFamilie_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	id_params, err := utils.ParamsChiRouter("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params of id!", err.Error())
		return
	}
	if id_params == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params of uuid!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.FamilieService.DeleteFamilie(ctx, id_params, middleware_token_id); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete the familie data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to delete the familie data!", true)

}

func (c *ControllerHandlerFamilie) UpdateFamilie_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.UpdateFamilie
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the familie payload!", err.Error())
		return
	}

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.FamilieService.UpdateFamilie(ctx, middleware_token_id, payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update the familie data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to update familie!", true)

}

func (c *ControllerHandlerFamilie) GetAllFamilie_Bp(w http.ResponseWriter, r *http.Request) {

	// Parse pagination from query params
	allowedSorts := []string{"name", "created_at"}
	params := utils.ParsePaginationParams(r, allowedSorts, "name")

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	paginatedData, err := c.FamilieService.GetAllFamilie(ctx, params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get families!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully retrieved families!", paginatedData)

}
