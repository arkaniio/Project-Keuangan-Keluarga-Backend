package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type GoalsService interface {
	CreateNewGoals(ctx context.Context, goals *model.Goals) error
	GetAllGoals(ctx context.Context, params model.PaginationParams, user_id uuid.UUID) (model.PaginatedResponse, error)
	DeleteGoals(ctx context.Context, user_id uuid.UUID) error
	UpdateGoals(ctx context.Context, user_id uuid.UUID, payload model.PayloadUpdateGoals) error
	TrackingProgressGoals(ctx context.Context, user_id uuid.UUID) ([]model.ProgressGoals, error)
	RemainingDaysGoals(ctx context.Context, user_id uuid.UUID) ([]model.RemainingDays, error)
}

type GoalsRepo struct {
	repo             repository.GoalsRepository
	repoUser         repository.UserRepository
	repoFamilyMember repository.FamilyMemberRepository
}

func NewGoalsService(repo repository.GoalsRepository, repoUser repository.UserRepository, repoFamilyMember repository.FamilyMemberRepository) GoalsService {
	return &GoalsRepo{repo: repo, repoUser: repoUser, repoFamilyMember: repoFamilyMember}
}

func (s *GoalsRepo) CreateNewGoals(ctx context.Context, goals *model.Goals) error {

	if goals.Current_amount >= goals.Target_amount {
		goals.Status = "completed"
	} else {
		goals.Status = "active"
	}

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, goals.User_id)
	if err != nil {
		return errors.New("User is not a member of any family")
	}
	goals.FamilyMemberId = fm.Id

	return s.repo.CreateNewGoals(ctx, goals)

}

func (s *GoalsRepo) GetAllGoals(ctx context.Context, params model.PaginationParams, user_id uuid.UUID) (model.PaginatedResponse, error) {

	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil {
		return model.PaginatedResponse{}, err
	}

	if fm == nil {
		return model.PaginatedResponse{
			Items: []model.PayloadGoalsWithUser{},
			Pagination: model.PaginationMeta{
				TotalItems:   0,
				TotalPages:   0,
				CurrentPage:  params.Page,
				PerPage:      params.Limit,
			},
		}, nil
	}

	goals_data, total_items, err := s.repo.GetAllGoals(ctx, params, fm.FamilyId)
	if err != nil {
		return model.PaginatedResponse{}, errors.New("Failed to get the all goals with the pagination")
	}

	meta := utils.BuildPaginationMeta(total_items, params.Page, params.Limit)

	return model.PaginatedResponse{
		Items:      goals_data,
		Pagination: meta,
	}, nil

}

func (s *GoalsRepo) DeleteGoals(ctx context.Context, user_id uuid.UUID) error {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != user_id {
		return errors.New("Failed to delete goals! because the id is not same!")
	}

	return s.repo.DeleteGoals(ctx, user_id)
}

func (s *GoalsRepo) UpdateGoals(ctx context.Context, user_id uuid.UUID, payload model.PayloadUpdateGoals) error {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return errors.New("Failed to get the users data in db!!")
	}

	if users_data.Id != user_id {
		return errors.New("Failed to update goals! because the id is not same!")
	}

	return s.repo.UpdateGoals(ctx, user_id, payload)
}

func (s *GoalsRepo) TrackingProgressGoals(ctx context.Context, user_id uuid.UUID) ([]model.ProgressGoals, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return []model.ProgressGoals{}, nil
	}
	return s.repo.TrackingProgressGoals(ctx, fm.FamilyId)
}

func (s *GoalsRepo) RemainingDaysGoals(ctx context.Context, user_id uuid.UUID) ([]model.RemainingDays, error) {
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, user_id)
	if err != nil || fm == nil {
		return []model.RemainingDays{}, nil
	}
	return s.repo.RemainingDaysGoals(ctx, fm.FamilyId)
}
