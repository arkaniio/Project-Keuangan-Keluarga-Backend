package controller

import (
	"context"
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

	if err := utils.ValidatePayloads(payloads); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate JSON", err.Error())
		return
	}

	users, err := utils.ParsingPayloadUser(payloads)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to exect the payload to user!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := s.service.CreateNewUser(ctx, users); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create new user", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to create new user", nil)

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

	users, err := s.service.GetUserByEmail(payloadsLogin.Email)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get user by email", err.Error())
		return
	}
	if users == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Email not found", nil)
		return
	}

	if err := utils.VerifyPassword(payloadsLogin.Password, users.Password); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to compare password", err.Error())
		return
	}

	token, err := utils.GenerateJwt(users.Id, users.Email, users.Username, users.Role)
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

func (s *ControllerHandler) GetUsersById(w http.ResponseWriter, r *http.Request) {

	id_params, err := utils.ParamsChiRouter("id", r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id params!", err.Error())
		return
	}
	if id_params == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get and detect the id params!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	users_data, err := s.service.GetUserById(ctx, id_params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the users data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to get the users data!", users_data)

}

func (s *ControllerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

	user_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
		return
	}
	if user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", false)
		return
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

	utils.ResponseSuccess(w, http.StatusOK, "Success to get the user id from token!", users)

}

func (s *ControllerHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	user_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user token from middleware!", err.Error())
		return
	}
	if user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid type, nil result!", false)
		return
	}

	if err := utils.ParsingMultipartFormData(w, r); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the multipart form data!", err.Error())
		return
	}

	var paylod model.UpdatePayloadUser
	value_payload, err := utils.ParsingFormValue(r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the form value!", err.Error())
		return
	}

	profile_img, header, err := r.FormFile("profile_img")
	if err != nil {
		if err == http.ErrMissingFile {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to detect file!", err.Error())
			return
		}
	}

	if err == nil {
		buff := make([]byte, 512)
		if _, err := profile_img.Read(buff); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to read file profile image!", err.Error())
			return
		}

		if err := utils.CheckRightPath(buff); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the content type!", err.Error())
			return
		}

		path, err := utils.MakeFileName("uploadsProfile", header, profile_img)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to make the file name!", err.Error())
			return
		}
		if path == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to make the file name!", false)
			return
		}

		ctx, cancle := context.WithTimeout(r.Context(), time.Second)
		defer cancle()

		user_data, err := s.service.GetUserById(ctx, user_id)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", err.Error())
			return
		}
		if user_data == nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token!", false)
			return
		}

		if err := utils.CheckOldPath(user_data.Profile_img); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to check the old path!", err.Error())
			return
		}

		utils.PayloaUpdate(&paylod.Profile_img, path)

	}

	utils.PayloaUpdate(&paylod.Username, value_payload.Name)
	utils.PayloaUpdate(&paylod.Email, value_payload.Email)
	utils.PayloaUpdate(&paylod.Password, value_payload.Password)
	utils.PayloaUpdate(&paylod.Role, value_payload.Role)

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := s.service.UpdateDataUser(user_id, ctx, paylod); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update the user data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to update the user data!", nil)

}

func (s *ControllerHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {

	middleware_role, err := middleware.GetTokenRole(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware role!", err.Error())
		return
	}
	if middleware_role == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the value of middleware role!", false)
		return
	}

	if middleware_role != "kepala keluarga" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this data!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	users, err := s.service.GetAllUser(ctx)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user data!", err.Error())
		return
	}
	if users == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user data!", false)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to get the user data!", users)

}
