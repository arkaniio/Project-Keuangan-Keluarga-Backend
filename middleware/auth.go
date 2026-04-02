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

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the strings token for auth!", false)
			return
		}

		claims_data, err := utils.ValidateToken(token)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the token!", err.Error())
			return
		}

		user_id_parse, err := uuid.Parse(claims_data.Id)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to parse the id!", err.Error())
			return
		}
		ctx_id := context.WithValue(r.Context(), "id", user_id_parse)
		r = r.WithContext(ctx_id)

		ctx_role := context.WithValue(r.Context(), "role", claims_data.Role)
		r = r.WithContext(ctx_role)

		next.ServeHTTP(w, r)

	})
}

func GetTokenId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

	data_id := r.Context().Value("id")
	user_id, ok := data_id.(uuid.UUID)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert into a uuid", false)
		return uuid.Nil, errors.New("Failed to convert into a uuid!")
	}
	if user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id from the token!", false)
		return uuid.Nil, errors.New("Failed to get the id from the token!")
	}

	return user_id, nil

}

func GetTokenRole(w http.ResponseWriter, r *http.Request) (string, error) {

	data_role := r.Context().Value("role")
	user_role, ok := data_role.(string)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "failed to convert into a string", false)
		return "", errors.New("Failed to convert into a string!")
	}
	if user_role == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role from the token!", false)
		return "", errors.New("Failed to get the role from the token!")
	}

	return user_role, nil

}
