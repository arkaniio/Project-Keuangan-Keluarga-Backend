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
	repo     repository.GoalsRepository
	repoUser repository.UserRepository
}

func NewGoalsService(repo repository.GoalsRepository, repoUser repository.UserRepository) GoalsService {
	return &GoalsRepo{repo: repo, repoUser: repoUser}
}

func (s *GoalsRepo) CreateNewGoals(ctx context.Context, goals *model.Goals) error {

	if goals.Current_amount >= goals.Target_amount {
		goals.Status = "completed"
	} else {
		goals.Status = "active"
	}

	return s.repo.CreateNewGoals(ctx, goals)

}

func (s *GoalsRepo) GetAllGoals(ctx context.Context, params model.PaginationParams, user_id uuid.UUID) (model.PaginatedResponse, error) {

	users_data, err := s.repoUser.GetUserById(ctx, user_id)
	if err != nil {
		return model.PaginatedResponse{}, errors.New("Failed to get the users data!")
	}

	if users_data.Role != "kepala keluarga" {
		return model.PaginatedResponse{}, errors.New("Failed to get the paginate response!")
	}

	goals_data, total_items, err := s.repo.GetAllGoals(ctx, params, user_id)
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
	return s.repo.TrackingProgressGoals(ctx, user_id)
}

func (s *GoalsRepo) RemainingDaysGoals(ctx context.Context, user_id uuid.UUID) ([]model.RemainingDays, error) {
	return s.repo.RemainingDaysGoals(ctx, user_id)
}
