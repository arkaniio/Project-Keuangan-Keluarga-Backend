package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewDecoder(r.Body).Decode(data)
}
