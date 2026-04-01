package utils

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ParamsMux(id string, r *http.Request) (uuid.UUID, error) {

	vars_id := mux.Vars(r)
	id_value := vars_id[id]

	uuid_final, err := uuid.Parse(id_value)
	if err != nil {
		return uuid.Nil, errors.New("Failed to convert data into an uuid!" + err.Error())
	}

	return uuid_final, nil

}
