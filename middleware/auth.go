package middleware

import (
	"context"
	"errors"
	"net/http"
	"project-keuangan-keluarga/utils"
	"strings"

	"github.com/google/uuid"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized", "Token is required")
			return
		}

		token := strings.TrimPrefix("Bearer ", header)
		if token == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the strings token for auth!", false)
			return
		}

		claims_data, err := utils.ValidateToken(token)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the token!", err.Error())
			return
		}

		ctx_id := context.WithValue(r.Context(), "id", claims_data.Id)
		ctx_role := context.WithValue(r.Context(), "role", claims_data.Role)
		r = r.WithContext(ctx_id)
		r = r.WithContext(ctx_role)

		next.ServeHTTP(w, r)

	})
}

func GetTokenId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

	data_id, ok := r.Context().Value("id").(uuid.UUID)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "failed to convert into a uuid", false)
		return uuid.Nil, errors.New("Failed to convert into a uuid!")
	}
	if data_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id from the token!", false)
		return uuid.Nil, errors.New("Failed to get the id from the token!")
	}

	return data_id, nil

}

func GetTokenRole(w http.ResponseWriter, r *http.Request) (string, error) {

	data_role, ok := r.Context().Value("role").(string)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "failed to convert into a string", false)
		return "", errors.New("Failed to convert into a string!")
	}
	if data_role == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role from the token!", false)
		return "", errors.New("Failed to get the role from the token!")
	}

	return data_role, nil

}
