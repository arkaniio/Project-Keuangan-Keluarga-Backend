package service

import (
	"context"
	"errors"
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

	users_data, err := s.repoUser.GetUserById(ctx, categories.UserId)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Role != "kepala keluarga" {
		return errors.New("Failed to access this method!")
	}

	return s.repoCategory.CreateNewCategory(ctx, categories)

}

func (s *repoCategoryCombine) UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory) error {

	users, err := s.repoUser.GetUserById(ctx, id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users.Role != "kepala keluarga" {
		return errors.New("Failed to access this method!")
	}

	return s.repoCategory.UpdateCategory(ctx, id, payload)

}

func (s *repoCategoryCombine) DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	users, err := s.repoUser.GetUserById(ctx, id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users.Role != "kepala keluarga" {
		return errors.New("Failed to access this method!")
	}

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
