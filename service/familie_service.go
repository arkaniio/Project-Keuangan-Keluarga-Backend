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
	repoFamilie      repository.FamilieRepository
	repoUser         repository.UserRepository
	repoFamilyMember repository.FamilyMemberRepository
	repoCategory     repository.CategoryRepository
}

func NewFamilieService(repoFamilie repository.FamilieRepository, repoUser repository.UserRepository, repoFamilyMember repository.FamilyMemberRepository, repoCategory repository.CategoryRepository) FamilieService {
	return &repoCombineFamilieAndUser{
		repoFamilie:      repoFamilie,
		repoUser:         repoUser,
		repoFamilyMember: repoFamilyMember,
		repoCategory:     repoCategory,
	}
}

func (s *repoCombineFamilieAndUser) CreateNewFamilie(ctx context.Context, familie *model.Familie) error {

	// Check if user is already a member of any family
	existingMember, _ := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, familie.Created_By)
	if existingMember != nil {
		return errors.New("User is already a member of a family. Cannot create a new one.")
	}

	users_data, err := s.repoUser.GetUserById(ctx, familie.Created_By)
	if err != nil {
		return errors.New("Failed to get the users data!")
	}

	if users_data.Id != familie.Created_By {
		return errors.New("Failed to access this method the id is not same!")
	}

	if err := s.repoFamilie.CreateNewFamilie(ctx, familie); err != nil {
		return err
	}

	// After successful family creation, seed default categories
	// We need the family_member_id of the creator (Kepala Keluarga)
	fm, err := s.repoFamilyMember.GetFamilyMemberByUserId(ctx, familie.Created_By)
	if err == nil && fm != nil {
		s.seedDefaultCategories(ctx, fm)
	}

	return nil

}

func (s *repoCombineFamilieAndUser) seedDefaultCategories(ctx context.Context, fm *model.FamilyMember) {
	defaults := []struct {
		Name string
		Type string
	}{
		{"Gaji", "income"},
		{"Bonus", "income"},
		{"Investasi", "income"},
		{"Makan & Minum", "expense"},
		{"Transportasi", "expense"},
		{"Belanja Rumah Tangga", "expense"},
		{"Pendidikan", "expense"},
		{"Kesehatan", "expense"},
		{"Cicilan", "expense"},
		{"Tagihan (Listrik/Air)", "expense"},
		{"Hiburan", "expense"},
		{"Lainnya", "expense"},
	}

	for _, d := range defaults {
		cat := &model.Category{
			Id:             uuid.New(),
			UserId:         fm.UserId,
			FamilyMemberId: fm.Id,
			Name:           d.Name,
			Type:           d.Type,
		}
		_ = s.repoCategory.CreateNewCategory(ctx, cat) // Best effort seeding
	}
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
