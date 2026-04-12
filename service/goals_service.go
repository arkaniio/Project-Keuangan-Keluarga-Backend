package service

import (
	"context"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type GoalsService interface {
	CreateNewGoals(ctx context.Context, goals *model.Goals) error
}

type GoalsRepo struct {
	repo repository.GoalsRepository
}

func NewGoalsService(repo repository.GoalsRepository) GoalsService {
	return &GoalsRepo{repo: repo}
}

func (s *GoalsRepo) CreateNewGoals(ctx context.Context, goals *model.Goals) error {
	return s.repo.CreateNewGoals(ctx, goals)
}
