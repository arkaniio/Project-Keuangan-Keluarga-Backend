package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

type ControllerHandler struct {
	service service.UserService
}

func NewUserController(svc service.UserService) *ControllerHandler {
	return &ControllerHandler{service: svc}
}

func (s *ControllerHandler) Register(w http.ResponseWriter, r *http.Request) {

	var payloads model.Payload
	if err := utils.DecodeJson(r, &payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode JSON", err.Error())
		return
	}

	if err := utils.IsValidEmail(payloads.Email); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate email", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate JSON", err.Error())
		return
	}

	users, err := utils.ParsingPayloadUser(payloads)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to exect the payload to user!", err.Error())
		return
	}

	users_email, err := s.service.GetUserByEmail(users.Email)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get user by email", err.Error())
		return
	}
	if users_email != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Email already exists", false)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to create new user", nil)

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := s.service.CreateNewUser(ctx, users); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create new user", err.Error())
		return
	}

}

func (s *ControllerHandler) Login(w http.ResponseWriter, r *http.Request) {

	var payloadsLogin model.LoginPayload
	if err := utils.DecodeJson(r, &payloadsLogin); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode JSON", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payloadsLogin); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate JSON", err.Error())
		return
	}

	if err := utils.IsValidEmail(payloadsLogin.Email); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate email", err.Error())
		return
	}

	users, err := s.service.GetUserByEmail(payloadsLogin.Email)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get user by email", err.Error())
		return
	}
	if users == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Email not found", false)
		return
	}

	if err := utils.VerifyPassword(payloadsLogin.Password, users.Password); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to compare password", err.Error())
		return
	}

	token, err := utils.GenerateJwt(users.Id, users.Email, users.Name, users.Role)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to generate token", err.Error())
		return
	}
	if token == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to generate token!", false)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to login", token)

}

func (s *ControllerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

	user_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", false)
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	users, err := s.service.GetUserById(ctx, user_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if users == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", false)
		return
	}

	//debug
	fmt.Println(users)

	utils.ResponseSuccess(w, http.StatusOK, "Success to get the user id from token!", users)

}
