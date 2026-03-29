package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DecodeJson(r *http.Request, v any) error {

	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type must be application/json")
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		return err
	}

	return nil
}
