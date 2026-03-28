package service

import (
	"context"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type UserService interface {
	CreateNewUser(ctx context.Context, user *model.User) error
}

type repoUser struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &repoUser{repo: repo}
}

// CreateNewUser delegates user creation to the repository layer.
func (s *repoUser) CreateNewUser(ctx context.Context, user *model.User) error {
	// TODO: add business logic / validation here
	return s.repo.CreateNewUser(ctx, user)
}
