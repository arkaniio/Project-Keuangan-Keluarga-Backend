package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FieldsFile struct {
	Name        string
	Email       string
	Password    string
	Role        string
	Profile_img string
}

func DetectContentType(buff []byte) string {

	content_type := http.DetectContentType(buff)
	if content_type != "jpg" && content_type != "jpeg" {
		return "Failed to detect the content type of the file!"
	}

	return ""

}

func CheckOldPath(value string) error {

	if value != "" {
		path_old := value
		if _, err := os.Stat(value); os.IsNotExist(err) {
			return errors.New("Failed to get the old path and check the old path is exist or not!")
		}
		if err := os.Remove(path_old); err != nil {
			return errors.New("Failed to remove the old path!")
		}
	}

	return nil

}

func MakeFileName(value string, form *multipart.FileHeader, file multipart.File) (string, error) {

	file_name := uuid.New().String() + form.Filename
	if err := os.MkdirAll(value, os.ModePerm); err != nil {
		return "", err
	}
	path := filepath.Join(value, file_name)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return path, nil

}

func CheckRightPath(filename []byte) error {

	check_file := http.DetectContentType(filename)

	if check_file != "jpg" && check_file != "png" && check_file != "jpeg" {
		return errors.New("Failed to detect the content type of the file!")
	}

	return nil

}

func ParsingMultipartFormData(w http.ResponseWriter, r *http.Request) error {

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return errors.New("Failed to get the body and parsing into an multipart form data!")
	}

	return nil

}

func ParsingFormValue(r *http.Request) (FieldsFile, error) {

	payloads_file_value := FieldsFile{
		Name:        r.FormValue("name"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		Role:        r.FormValue("role"),
		Profile_img: r.FormValue("profile_img"),
	}

	return payloads_file_value, nil

}
