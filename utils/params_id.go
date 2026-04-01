package utils

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func ParamsChiRouter(id string, r *http.Request) (uuid.UUID, error) {

	url := chi.URLParam(r, id)
	id_parse, err := uuid.Parse(url)
	if err != nil {
		return uuid.Nil, errors.New("Failed to get the uuid parse!" + err.Error())
	}

	return id_parse, nil

}
