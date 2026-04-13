package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"
)

type GoalsService interface {
	CreateNewGoals(ctx context.Context, goals *model.Goals) error
	GetAllGoals(ctx context.Context, params model.PaginationParams) (model.PaginatedResponse, error)
}

type GoalsRepo struct {
	repo repository.GoalsRepository
}

func NewGoalsService(repo repository.GoalsRepository) GoalsService {
	return &GoalsRepo{repo: repo}
}

func (s *GoalsRepo) CreateNewGoals(ctx context.Context, goals *model.Goals) error {

	if goals.Current_amount >= goals.Target_amount {
		goals.Status = "completed"
	} else {
		goals.Status = "active"
	}

	return s.repo.CreateNewGoals(ctx, goals)

}

func (s *GoalsRepo) GetAllGoals(ctx context.Context, params model.PaginationParams) (model.PaginatedResponse, error) {

	goals_data, total_items, err := s.repo.GetAllGoals(ctx, params)
	if err != nil {
		return model.PaginatedResponse{}, errors.New("Failed to get the all goals with the pagination")
	}

	meta := utils.BuildPaginationMeta(total_items, params.Page, params.Limit)

	return model.PaginatedResponse{
		Items:      goals_data,
		Pagination: meta,
	}, nil

}
