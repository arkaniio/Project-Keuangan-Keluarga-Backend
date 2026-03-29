package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type KeuanganService interface {
	CreateNewKeuangan(ctx context.Context, keuangan *model.Keuangan) error
}

type repoKeuangan struct {
	repo repository.KeuanganRepository
}

func NewKeuanganService(repo repository.KeuanganRepository) KeuanganService {
	return &repoKeuangan{repo: repo}
}

func (s *repoKeuangan) CreateNewKeuangan(ctx context.Context, keuangan *model.Keuangan) error {
	return s.repo.CreateNewKeuangan(ctx, keuangan)
}
