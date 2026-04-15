package service

import (
	"context"
	"errors"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/repository"
)

type FamilieService interface {
	CreateNewFamilie(ctx context.Context, familie *model.Familie) error
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
