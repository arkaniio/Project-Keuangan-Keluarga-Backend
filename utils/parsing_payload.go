package utils

import (
	"project-keuangan-keluarga/model"
	"time"

	"github.com/google/uuid"
)

func ParsingPayloadUser(payload model.Payload) (*model.User, error) {

	hashing_password, err := HashPassword(payload.Password)
	if err != nil {
		return &model.User{}, err
	}

	return &model.User{
		Id:          uuid.New(),
		Name:        payload.Name,
		Email:       payload.Email,
		Password:    hashing_password,
		Role:        payload.Role,
		Profile_img: payload.Profile_img,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func ParsingPayloadTransaction(payload model.PayloadTransaction) (*model.Transaction, error) {

	return &model.Transaction{
		Id:          uuid.New(),
		UserId:      payload.UserId,
		Type:        payload.Type,
		Amount:      payload.Amount,
		CategoryId:  payload.CategoryId,
		Description: payload.Description,
		Date:        payload.Date,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func ParsingPayloadCategory(payload model.PayloadCategory) (*model.Category, error) {

	return &model.Category{
		Id:     uuid.New(),
		UserId: payload.UserId,
		Name:   payload.Name,
		Type:   payload.Type,
	}, nil
}

func PayloaUpdate(dest **string, val string) {

	if val != "" {
		*dest = &val
	}

}

func PayloaUpdateInt64(dest **int64, val int64) {

	if val != 0 {
		*dest = &val
	}

}
