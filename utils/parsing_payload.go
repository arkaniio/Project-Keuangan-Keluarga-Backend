package utils

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"project-keuangan-keluarga/model"
)

func ParsingPayloadUser(payload model.Payload) (*model.User, error) {

	hashing_password, err := HashPassword(payload.Password)
	if err != nil {
		return &model.User{}, err
	}

	return &model.User{
		Id:          uuid.New(),
		Username:    payload.Username,
		Email:       payload.Email,
		Password:    hashing_password,
		Role:        payload.Role,
		Profile_img: payload.Profile_img,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func ParsingPayloadTransaction(payload model.PayloadTransaction, userId uuid.UUID) (*model.Transaction, error) {

	if userId == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.Transaction{
		Id:          uuid.New(),
		UserId:      userId,
		Type:        payload.Type,
		Amount:      payload.Amount,
		CategoryId:  payload.CategoryId,
		Description: payload.Description,
		Date:        payload.Date,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func ParsingPayloadBudget(payload model.PayloadBudget, userId uuid.UUID) (*model.Budget, error) {

	if userId == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.Budget{
		Id:           uuid.New(),
		UserId:       userId,
		Category_Id:  payload.Category_Id,
		Limit_amount: payload.Limit_amount,
		Period:       payload.Period,
		StartDate:    payload.StartDate,
		EndDate:      payload.EndDate,
		IsActive:     payload.IsActive,
		Created_at:   time.Now().UTC(),
		Updated_at:   time.Now().UTC(),
	}, nil

}

func ParsingPayloadCategory(payload model.PayloadCategory, userId uuid.UUID) (*model.Category, error) {

	if userId == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.Category{
		Id:     uuid.New(),
		UserId: userId,
		Name:   payload.Name,
		Type:   payload.Type,
	}, nil
}

func PayloaUpdate(dest **string, val string) {

	if val != "" {
		*dest = &val
	}

}

func PayloadJoinDataTransaction(payload model.PayloadTransactionDataCategory) (*model.PayloadTransactionWithCategory, error) {

	if payload.Id == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.PayloadTransactionWithCategory{
		Id:          payload.Id,
		UserId:      payload.UserId,
		Type:        payload.Type,
		Amount:      payload.Amount,
		CategoryId:  payload.CategoryId,
		Description: payload.Description,
		Date:        payload.Date,
		CreatedAt:   payload.CreatedAt,
		UpdatedAt:   payload.UpdatedAt,
		Category: model.Category{
			Name: payload.Name,
		},
	}, nil

}

func PayloadJoinDataCategory(payload model.PayloadCategoryWithUserData) (model.PayloadCategoryWithUser, error) {

	if payload.Id == uuid.Nil {
		return model.PayloadCategoryWithUser{}, errors.New("Failed to get the uuid file type!")
	}

	return model.PayloadCategoryWithUser{
		Id:     payload.Id,
		UserId: payload.UserId,
		User: model.User{
			Id:       payload.UserId,
			Username: payload.Username,
			Email:    payload.Email,
		},
		Name: payload.CategoryName,
		Type: payload.CategoryType,
	}, nil

}

func PayloadJoinDataCategoryAndUser(payload model.BudgetWithCategoryAndUser) (*model.BudgetWithCategoryAndUserData, error) {

	if payload.Id == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.BudgetWithCategoryAndUserData{
		Id:     payload.Id,
		UserId: payload.UserId,
		User_Data: model.User{
			Username: payload.User_Name,
			Email:    payload.User_Email,
		},
		Category_Data: model.Category{
			Name: payload.Category_Name,
			Type: payload.Category_Type,
		},
		Limit_amount: payload.Limit_amount,
		Period:       payload.Period,
		StartDate:    payload.StartDate,
		EndDate:      payload.EndDate,
		IsActive:     payload.IsActive,
	}, nil

}

func ParsingPayloadGoals(payload model.PayloadGoals, userId uuid.UUID) (*model.Goals, error) {

	if userId == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.Goals{
		Id:             uuid.New(),
		User_id:        userId,
		Name:           payload.Name,
		Target_amount:  payload.Target_amount,
		Current_amount: payload.Current_amount,
		Start_date:     payload.Start_date,
		Target_date:    payload.Target_date,
		Status:         payload.Status,
		Created_at:     time.Now().UTC(),
		Updated_at:     time.Now().UTC(),
	}, nil

}

func PayloadJoinDataGoals(payload model.PayloadGoalsWithUserData) (*model.PayloadGoalsWithUser, error) {

	if payload.Id == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.PayloadGoalsWithUser{
		Id:      payload.Id,
		User_id: payload.User_id,
		User: model.User{
			Username: payload.Username,
			Email:    payload.Email,
		},
		Name:           payload.Name,
		Target_amount:  payload.Target_amount,
		Current_amount: payload.Current_amount,
		Start_date:     payload.Start_date,
		Target_date:    payload.Target_date,
		Status:         payload.Status,
		Created_at:     payload.Created_at,
		Updated_at:     payload.Updated_at,
	}, nil

}

func ParsingPayloadFamilie(payload model.PayloadFamilie, userId uuid.UUID) (*model.Familie, error) {

	if userId == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.Familie{
		Id:         uuid.New(),
		Name:       payload.Name,
		Created_By: userId,
		Created_at: time.Now().UTC(),
	}, nil

}

func PayloadJoinDataFamilie(payload model.PayloadFamilieWithUserData) (model.PayloadFamilieWithUser, error) {

	if payload.Id == uuid.Nil {
		return model.PayloadFamilieWithUser{}, errors.New("Failed to get the uuid file type!")
	}

	return model.PayloadFamilieWithUser{
		Id:         payload.Id,
		Name:       payload.Name,
		Created_By: payload.Created_By,
		User: model.User{
			Id:       payload.Created_By,
			Username: payload.Username,
			Email:    payload.Email,
		},
		Created_at: payload.Created_at,
	}, nil

}

func PayloaUpdateInt64(dest **int64, val int64) {

	if val != 0 {
		*dest = &val
	}

}

func ParsingPayloadFamilyMember(payload model.PayloadFamilyMember, userId uuid.UUID) (*model.FamilyMember, error) {

	if userId == uuid.Nil {
		return nil, errors.New("Failed to get the uuid file type!")
	}

	return &model.FamilyMember{
		Id:       uuid.New(),
		FamilyId: payload.FamilyId,
		UserId:   userId,
		Role:     payload.Role,
		JoinedAt: time.Now().UTC(),
	}, nil

}

func PayloadJoinDataFamilyMember(payload model.PayloadFamilyMemberWithUserData) (model.PayloadFamilyMemberWithUser, error) {

	if payload.Id == uuid.Nil {
		return model.PayloadFamilyMemberWithUser{}, errors.New("Failed to get the uuid file type!")
	}

	return model.PayloadFamilyMemberWithUser{
		Id:       payload.Id,
		FamilyId: payload.FamilyId,
		UserId:   payload.UserId,
		User: model.User{
			Id:       payload.UserId,
			Username: payload.Username,
			Email:    payload.Email,
		},
		Role:     payload.Role,
		JoinedAt: payload.JoinedAt,
	}, nil

}
