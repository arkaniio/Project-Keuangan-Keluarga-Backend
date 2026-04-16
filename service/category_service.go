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
	GetAllCategory(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
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
		return errors.New("User is not a member of any family")
	}

	if fm.Role != "kepala keluarga" {
		return errors.New("Unauthorized: Only kepala keluarga can create categories")
	}

	categories.FamilyMemberId = fm.Id

	return s.repoCategory.CreateNewCategory(ctx, categories)

}

func (s *repoCategoryCombine) UpdateCategory(ctx context.Context, id uuid.UUID, payload model.UpdatePayloadCategory, user_id uuid.UUID) error {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return errors.New("User is not a member of any family")
	}

	if fm.Role != "kepala keluarga" {
		return errors.New("Unauthorized: Only kepala keluarga can update categories")
	}

	return s.repoCategory.UpdateCategory(ctx, id, payload)

}

func (s *repoCategoryCombine) DeleteCategory(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return errors.New("User is not a member of any family")
	}

	if fm.Role != "kepala keluarga" {
		return errors.New("Unauthorized: Only kepala keluarga can delete categories")
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
