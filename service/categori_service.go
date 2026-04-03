package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateNewCategory(ctx context.Context, categories *model.Category) error
	UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error
}

type repoCategory struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &repoCategory{repo: repo}
}

func (s *repoCategory) CreateNewCategory(ctx context.Context, categories *model.Category) error {
	return s.repo.CreateNewCategory(ctx, categories)
}

func (s *repoCategory) UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error {
	return s.repo.UpdateCategory(ctx, id, payload)
}
