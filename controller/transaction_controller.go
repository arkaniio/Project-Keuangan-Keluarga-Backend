package controller

import (
	"project-keuangan-keluarga/service"
)

type ControllerHandlerKeuangan struct {
	KeuanganService service.KeuanganService
}

func NewControllerHandlerKeuangan(keuanganService service.KeuanganService) *ControllerHandlerKeuangan {
	return &ControllerHandlerKeuangan{KeuanganService: keuanganService}
}
