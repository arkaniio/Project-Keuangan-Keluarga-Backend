package service

import (
	"project-keuangan-keluarga/repository"
)

type KeuanganService interface {
}

type repoKeuangan struct {
	repo repository.KeuanganRepository
}

func NewKeuanganService(repo repository.KeuanganRepository) KeuanganService {
	return &repoKeuangan{repo: repo}
}
