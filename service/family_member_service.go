package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type FamilyMemberService interface {
	CreateFamilyMember(ctx context.Context, member *model.FamilyMember) error
	UpdateFamilyMember(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilyMember) error
	DeleteFamilyMember(ctx context.Context, user_id uuid.UUID) error
	GetAllFamilyMember(ctx context.Context, user_id uuid.UUID, params model.PaginationParams) (*model.PaginatedResponse, error)
	GetFamilyMemberByUserId(ctx context.Context, user_id uuid.UUID) (*model.FamilyMember, error)
}

type repoCombineFamilyMemberAndUser struct {
	repoFamilyMember repository.FamilyMemberRepository
	repoUser         repository.UserRepository
}

func NewFamilyMemberService(repoFamilyMember repository.FamilyMemberRepository, repoUser repository.UserRepository) FamilyMemberService {
	return &repoCombineFamilyMemberAndUser{
		repoFamilyMember: repoFamilyMember,
		repoUser:         repoUser,
	}
}

func (s *repoCombineFamilyMemberAndUser) CreateFamilyMember(ctx context.Context, member *model.FamilyMember) error {

	// Check if user is already a member of any family
	existingMember, _ := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, member.UserId)
	if existingMember != nil {
		return errors.New("User is already a member of a family. Cannot join another one.")
	}

	users_data, err := s.repoUser.GetUserById(ctx, member.UserId)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != member.UserId {
		return errors.New("Failed to access this method the id is not same!")
	}

	// Default role to "anggota" if not provided or to ensure non-kepala-keluarga joining
	if member.Role == "" || member.Role == "kepala keluarga" {
		member.Role = "anggota"
	}

	return s.repoFamilyMember.CreateFamilyMember(ctx, member)

}

func (s *repoCombineFamilyMemberAndUser) UpdateFamilyMember(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilyMember) error {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != user_id {
		return errors.New("Failed to access this method!")
	}

	return s.repoFamilyMember.UpdateFamilyMember(ctx, user_id, payload)

}

func (s *repoCombineFamilyMemberAndUser) DeleteFamilyMember(ctx context.Context, user_id uuid.UUID) error {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != user_id {
		return errors.New("Failed to access this method!")
	}

	return s.repoFamilyMember.DeleteFamilyMember(ctx, user_id)

}

func (s *repoCombineFamilyMemberAndUser) GetAllFamilyMember(ctx context.Context, user_id uuid.UUID, params model.PaginationParams) (*model.PaginatedResponse, error) {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return nil, err
	}

	if fm == nil {
		return &model.PaginatedResponse{
			Items: []model.PayloadFamilyMemberWithUser{},
			Pagination: model.PaginationMeta{
				TotalItems:   0,
				TotalPages:   0,
				CurrentPage:  params.Page,
				PerPage:      params.Limit,
			},
		}, nil
	}

	items, totalItems, err := s.repoFamilyMember.GetAllFamilyMember(ctx, fm.FamilyId, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}

func (s *repoCombineFamilyMemberAndUser) GetFamilyMemberByUserId(ctx context.Context, user_id uuid.UUID) (*model.FamilyMember, error) {
	return s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
}

