package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateNewCategory(ctx context.Context, categories *model.Category) error
	UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	GetAllCategory(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
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

func (s *repoCategory) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *repoCategory) GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	return s.repo.GetCategoryById(ctx, id)
}

func (s *repoCategory) GetAllCategory(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error) {
	items, totalItems, err := s.repo.GetAllCategory(ctx, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}
