package utils

import (
	"project-keuangan-keluarga/model"

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
		CreatedAt:   payload.CreatedAt,
		UpdatedAt:   payload.UpdatedAt,
	}, nil
}
