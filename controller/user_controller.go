package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	r.Body = http.MaxBytesReader(w, r.Body, 20>>10)

	if err := r.ParseMultipartForm(20 >> 10); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing into a multipart form data!", err.Error())
		return
	}

	var paylod model.PayloadUpdate
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

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
		profile_img.Seek(0, 0)

		content_type := http.DetectContentType(buff)
		if content_type != "jpg" && content_type != "jpeg" && content_type != "png" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the content type of the file!", "")
			return
		}

		file_name := uuid.New().String() + header.Filename
		path_folder := "/uploadsProfile"
		os.MkdirAll(path_folder, os.ModePerm)
		path := filepath.Join(path_folder, file_name)
		if path == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the path from os!", false)
			return
		}

		dst, err := os.Create(path)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to create the file!", err.Error())
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, profile_img); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to copy the file!", err.Error())
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

		if user_data.Profile_img != "" {
			path_old := user_data.Profile_img
			if _, err := os.Stat(path_old); os.IsNotExist(err) {
				utils.ResponseError(w, http.StatusBadRequest, "Failed to get the old profile image from the form data!", err.Error())
				return
			}
			if err := os.Remove(path_old); err != nil {
				utils.ResponseError(w, http.StatusBadRequest, "Failed to remove the old profile image from the form data!", err.Error())
				return
			}
		}

		utils.PayloaUpdate(&paylod.Profile_img, path)

	}

	utils.PayloaUpdate(&paylod.Name, name)
	utils.PayloaUpdate(&paylod.Email, email)
	utils.PayloaUpdate(&paylod.Password, password)
	utils.PayloaUpdate(&paylod.Role, role)

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := s.service.UpdateDataUser(user_id, ctx, paylod); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update the user data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success to update the user data!", nil)

}
