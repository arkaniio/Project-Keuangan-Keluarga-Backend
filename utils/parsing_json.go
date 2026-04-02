package utils

import (
	"encoding/json"
	"net/http"
)

//func to parsing json response

type JsonType struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseJson(w http.ResponseWriter, statusCode int, message string, data interface{}) JsonType {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JsonType{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
	return JsonType{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}
}

func ResponseSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) JsonType {
	return ResponseJson(w, statusCode, message, data)
}

func ResponseError(w http.ResponseWriter, statusCode int, message string, data interface{}) JsonType {
	return ResponseJson(w, statusCode, message, data)
}
