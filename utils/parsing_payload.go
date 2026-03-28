package utils

import "project-keuangan-keluarga/model"

func ParsingPayloadUser(payload model.Payload) (model.User, error) {
	return model.User{
		Id:        payload.Id,
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  payload.Password,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}, nil
}
