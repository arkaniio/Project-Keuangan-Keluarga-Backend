package controller

import (
	"net/http"

	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"

	"github.com/go-playground/validator/v10"
)

type ControllerHandler struct {
	service service.UserService
}

func NewUserController(svc service.UserService) *ControllerHandler {
	return &ControllerHandler{service: svc}
}

func (s *ControllerHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {

	var payloads model.Payload
	if err := utils.DecodeJson(w, r, &payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode JSON", err.Error())
		return
	}

	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(payloads); err != nil {
		var errors []string
		for _, Err := range err.(validator.ValidationErrors) {
			errors = append(errors, Err.Error())
		}
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate JSON", err.Error())
		return
	}

	users, err := utils.ParsingPayloadUser(payloads)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to exect the payload to user!", err.Error())
		return
	}

	if err := s.service.CreateNewUser(r.Context(), &users); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create new user", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to create new user", nil)

}
