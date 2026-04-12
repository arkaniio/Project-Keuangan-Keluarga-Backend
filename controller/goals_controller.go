package controller

import (
	"net/http"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
)

type ControllerGoals struct {
	service service.GoalsService
}

func NewControllerGoals(service *service.GoalsService) ControllerGoals {
	return ControllerGoals{service: *service}
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

}
