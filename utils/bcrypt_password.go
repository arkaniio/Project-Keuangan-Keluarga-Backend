package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hash_password), nil
}

func VerifyPassword(password string, hash_password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(password))
	if err != nil {
		return errors.New("failed to verify password!")
	}
	return nil
}
