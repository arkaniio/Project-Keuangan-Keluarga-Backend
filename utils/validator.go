package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidatePayloads(payload interface{}) error {
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		var Error []string
		for _, Err := range err.(validator.ValidationErrors) {
			Error = append(Error, fmt.Sprintf("Failed to validate the payloads! %v", Err.Error()))
		}
		return errors.New("Failed to setting the validator" + err.Error())
	}

	return nil
}
