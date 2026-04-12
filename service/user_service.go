package service

import (
	"context"
	"errors"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type UserService interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateDataUser(id uuid.UUID, ctx context.Context, user model.UpdatePayloadUser) error
	GetAllUser(ctx context.Context) ([]model.User, error)
}

type repoUser struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &repoUser{repo: repo}
}

func (s *repoUser) CreateNewUser(ctx context.Context, user *model.User) error {

	users_data, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return errors.New("Failed to get the users data based on their email!")
	}
	if users_data != nil {
		return errors.New("Failed to get the users data because the value of users data is nil!")
	}

	if err := utils.IsValidEmail(user.Email); err != nil {
		return errors.New("Failed to detect the right format of email user!")
	}

	return s.repo.CreateNewUser(ctx, user)

}

func (s *repoUser) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *repoUser) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *repoUser) UpdateDataUser(id uuid.UUID, ctx context.Context, payload model.UpdatePayloadUser) error {
	return s.repo.UpdateDataUser(id, ctx, payload)
}

func (s *repoUser) GetAllUser(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAllUser(ctx)
}
