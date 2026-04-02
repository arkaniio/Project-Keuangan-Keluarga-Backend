package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type CategoryService interface {
	CreateNewCategory(ctx context.Context, categories *model.Category) error
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
