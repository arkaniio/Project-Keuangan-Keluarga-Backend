package service

import (
	"context"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type UserService interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type repoUser struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &repoUser{repo: repo}
}

func (s *repoUser) CreateNewUser(ctx context.Context, user *model.User) error {
	return s.repo.CreateNewUser(ctx, user)
}

func (s *repoUser) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
