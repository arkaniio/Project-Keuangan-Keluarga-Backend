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
	UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory, user_id uuid.UUID) error
	DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	GetAllCategory(ctx context.Context, user_id uuid.UUID, params model.PaginationParams) (*model.PaginatedResponse, error)
}

type repoCategoryCombine struct {
	repoCategory     repository.CategoryRepository
	repoUser         repository.UserRepository
	repoFamilyMember repository.FamilyMemberRepository
}

func NewCategoryService(repoCategory repository.CategoryRepository, repoUser repository.UserRepository, repoFamilyMember repository.FamilyMemberRepository) CategoryService {
	return &repoCategoryCombine{repoCategory: repoCategory, repoUser: repoUser, repoFamilyMember: repoFamilyMember}
}

func (s *repoCategoryCombine) CreateNewCategory(ctx context.Context, categories *model.Category) error {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, categories.UserId)
	if err != nil {
		return errors.New("Failed to verify family membership")
	}

	if fm == nil {
		return errors.New("Unauthorized: User must be in a family to create categories")
	}

	categories.FamilyMemberId = fm.Id

	return s.repoCategory.CreateNewCategory(ctx, categories)

}

func (s *repoCategoryCombine) UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory, user_id uuid.UUID) error {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return errors.New("Failed to verify family membership")
	}

	if fm == nil {
		return errors.New("Unauthorized: User must be in a family to update categories")
	}

	return s.repoCategory.UpdateCategory(ctx, id, payload)

}

func (s *repoCategoryCombine) DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return errors.New("Failed to verify family membership")
	}

	if fm == nil {
		return errors.New("Unauthorized: User must be in a family to delete categories")
	}

	return s.repoCategory.DeleteCategory(ctx, id, user_id)

}

func (s *repoCategoryCombine) GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	return s.repoCategory.GetCategoryById(ctx, id)
}

func (s *repoCategoryCombine) GetAllCategory(ctx context.Context, user_id uuid.UUID, params model.PaginationParams) (*model.PaginatedResponse, error) {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return nil, err
	}

	if fm == nil {
		return &model.PaginatedResponse{
			Items: []model.PayloadCategoryWithUser{},
			Pagination: model.PaginationMeta{
				TotalItems:   0,
				TotalPages:   0,
				CurrentPage:  params.Page,
				PerPage:      params.Limit,
			},
		}, nil
	}

	items, totalItems, err := s.repoCategory.GetAllCategory(ctx, fm.FamilyId, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}
