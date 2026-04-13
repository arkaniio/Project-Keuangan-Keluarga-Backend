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

type ControllerGoals struct {
	service service.GoalsService
}

func NewControllerGoals(service service.GoalsService) ControllerGoals {
	return ControllerGoals{service: service}
}

func (c *ControllerGoals) CreateNewGoals_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.PayloadGoals
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode json the payload!", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payloads!", err.Error())
		return
	}

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the token id from middleware!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid type for middleware token id!", false)
	}

	goals, err := utils.ParsingPayloadGoals(payload, middleware_token_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing into an goals model db!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.service.CreateNewGoals(ctx, goals); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create new goals for a user!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to create the new goals!", true)

}

func (c *ControllerGoals) GetAllGoals_Bp(w http.ResponseWriter, r *http.Request) {

	allowed_sort := []string{"created_at", "target_amount", "current_amount"}

	parsing_params := utils.ParsePaginationParams(r, allowed_sort, "created_at")

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	goals_data, err := c.service.GetAllGoals(ctx, parsing_params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the data of goals!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to get the full data of goals", goals_data)

}

func (c *ControllerGoals) DeleteGoals_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware id token!", err.Error())
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.service.DeleteGoals(ctx, middleware_token_id); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete the goals!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to delete the goals!", true)

}
