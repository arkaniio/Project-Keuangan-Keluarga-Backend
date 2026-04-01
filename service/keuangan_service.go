package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"

	"github.com/google/uuid"
)

type KeuanganService interface {
	CreateNewKeuangan(ctx context.Context, keuangan *model.Keuangan) error
	DeleteDataKeuangan(ctx context.Context, id uuid.UUID) error
	UpdateKeuangan(ctx context.Context, id uuid.UUID, payload model.PaylodUpdateKeuangan) error
	GetAllKeuangans(ctx context.Context) ([]model.KeuanganDataWithUser, error)
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

func (s *repoKeuangan) DeleteDataKeuangan(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDataKeuangan(ctx, id)
}

func (s *repoKeuangan) UpdateKeuangan(ctx context.Context, id uuid.UUID, payload model.PaylodUpdateKeuangan) error {
	return s.repo.UpdateKeuangan(ctx, id, payload)
}

func (s *repoKeuangan) GetAllKeuangans(ctx context.Context) ([]model.KeuanganDataWithUser, error) {
	return s.repo.GetAllKeuangans(ctx)
}
