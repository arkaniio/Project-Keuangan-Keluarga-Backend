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
	DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	GetAllCategory(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
}

type repoCategoryCombine struct {
	repoCategory repository.CategoryRepository
	repoUser     repository.UserRepository
}

func NewCategoryService(repoCategory repository.CategoryRepository, repoUser repository.UserRepository) CategoryService {
	return &repoCategoryCombine{repoCategory: repoCategory, repoUser: repoUser}
}

func (s *repoCategoryCombine) CreateNewCategory(ctx context.Context, categories *model.Category) error {

	return s.repoCategory.CreateNewCategory(ctx, categories)

}

func (s *repoCategoryCombine) UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error {

	return s.repoCategory.UpdateCategory(ctx, id, payload)

}

func (s *repoCategoryCombine) DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	return s.repoCategory.DeleteCategory(ctx, id, user_id)

}

func (s *repoCategoryCombine) GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	return s.repoCategory.GetCategoryById(ctx, id)
}

func (s *repoCategoryCombine) GetAllCategory(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error) {
	items, totalItems, err := s.repoCategory.GetAllCategory(ctx, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}
