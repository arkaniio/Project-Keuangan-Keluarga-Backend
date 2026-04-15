package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type FamilieService interface {
	CreateNewFamilie(ctx context.Context, familie *model.Familie) error
	DeleteFamilie(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error
	UpdateFamilie(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilie) error
	GetAllFamilie(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error)
}

type repoCombineFamilieAndUser struct {
	repoFamilie repository.FamilieRepository
	repoUser    repository.UserRepository
}

func NewFamilieService(repoFamilie repository.FamilieRepository, repoUser repository.UserRepository) FamilieService {
	return &repoCombineFamilieAndUser{
		repoFamilie: repoFamilie,
		repoUser:    repoUser,
	}
}

func (s *repoCombineFamilieAndUser) CreateNewFamilie(ctx context.Context, familie *model.Familie) error {

	users_data, err := s.repoUser.GetUserById(ctx, familie.Created_By)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != familie.Created_By {
		return errors.New("Failed to access this method the id is not same!")
	}

	return s.repoFamilie.CreateNewFamilie(ctx, familie)

}

func (s *repoCombineFamilieAndUser) DeleteFamilie(ctx context.Context, id uuid.UUID, user_id uuid.UUID) error {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != user_id {
		return errors.New("Failed to access this method!")
	}

	return s.repoFamilie.DeleteFamilie(ctx, id, user_id)

}

func (s *repoCombineFamilieAndUser) UpdateFamilie(ctx context.Context, user_id uuid.UUID, payload model.UpdateFamilie) error {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != user_id {
		return errors.New("Failed to access this method!")
	}

	return s.repoFamilie.UpdateFamilie(ctx, user_id, payload)

}

func (s *repoCombineFamilieAndUser) GetAllFamilie(ctx context.Context, params model.PaginationParams) (*model.PaginatedResponse, error) {
	items, totalItems, err := s.repoFamilie.GetAllFamilie(ctx, params)
	if err != nil {
		return nil, err
	}

	meta := utils.BuildPaginationMeta(totalItems, params.Page, params.Limit)

	return &model.PaginatedResponse{
		Items:      items,
		Pagination: meta,
	}, nil
}
