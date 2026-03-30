package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ParamsMux(id string, r *http.Request) (uuid.UUID, error) {

	vars_id := mux.Vars(r)
	user_id := vars_id["user_id"]

	return uuid.Parse(user_id)

}
