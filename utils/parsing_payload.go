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

func PayloaUpdateInt64(dest **int64, val int64) {

	if val != 0 {
		*dest = &val
	}

}
