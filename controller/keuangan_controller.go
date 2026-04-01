package controller

import (
	"context"
	"net/http"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
	"time"
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

func (s *ControllerHandlerKeuangan) DeleteKeuangan(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParamsMux("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parse the id", err.Error())
		return
	}

	if err := s.KeuanganService.DeleteDataKeuangan(r.Context(), id); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to delete the data", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to delete the data", nil)

}

func (s *ControllerHandlerKeuangan) UpdateKeuanganData(w http.ResponseWriter, r *http.Request) {

	var payloads model.PaylodUpdateKeuangan
	var keuangans model.Keuangan
	if err := utils.DecodeJson(r, &payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the json type!", err.Error())
		return
	}

	keuangan_id, err := utils.ParamsMux("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to setting the mux params for this method!", err.Error())
		return
	}

	utils.PayloaUpdate(&payloads.JenisTransaksi, keuangans.JenisTransaksi)
	utils.PayloaUpdateInt64(&payloads.JumlahPengeluaran, keuangans.JumlahPengeluaran)
	utils.PayloaUpdateInt64(&payloads.JumlahPemasukan, keuangans.JumlahPemasukan)
	utils.PayloaUpdate(&payloads.Kategori, keuangans.Kategori)
	utils.PayloaUpdate(&payloads.Tanggal, keuangans.Tanggal)

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := s.KeuanganService.UpdateKeuangan(ctx, keuangan_id, payloads); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update the data", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to update the data", nil)

}
