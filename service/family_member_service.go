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
	GetAllFamilyMember(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
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

	users_data, err := s.repoUser.GetUserById(ctx, member.UserId)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != member.UserId {
		return errors.New("Failed to access this method the id is not same!")
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

func (s *repoCombineFamilyMemberAndUser) GetAllFamilyMember(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error) {
	items, totalItems, err := s.repoFamilyMember.GetAllFamilyMember(ctx, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}
