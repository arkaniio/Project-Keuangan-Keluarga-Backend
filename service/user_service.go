package service

import (
	"context"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"

	"github.com/google/uuid"
)

type UserService interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateDataUser(id uuid.UUID, ctx context.Context, user model.PayloadUpdate) error
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

func (s *repoUser) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *repoUser) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *repoUser) UpdateDataUser(id uuid.UUID, ctx context.Context, user model.PayloadUpdate) error {
	return s.repo.UpdateDataUser(id, ctx, user)
}
