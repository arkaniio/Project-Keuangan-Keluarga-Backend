package controller

import (
	"net/http"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
)

type ControllerHandlerKeuangan struct {
	KeuanganService service.KeuanganService
}

func NewControllerHandlerKeuangan(keuanganService service.KeuanganService) *ControllerHandlerKeuangan {
	return &ControllerHandlerKeuangan{KeuanganService: keuanganService}
}

func (c *ControllerHandlerKeuangan) CreateNewKeuangan(w http.ResponseWriter, r *http.Request) {

	var payloads model.PayloadKeuangan
	if err := utils.DecodeJson(r, &payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode JSON", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate JSON", err.Error())
		return
	}

	keuangan, err := utils.ParsingPayloadKeuangan(payloads)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the payloads to keuangan!", err.Error())
		return
	}

	if err := c.KeuanganService.CreateNewKeuangan(r.Context(), keuangan); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create new keuangan!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to create new keuangan", nil)

}
