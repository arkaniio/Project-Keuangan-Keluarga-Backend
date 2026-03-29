package utils

import "net/mail"

func IsValidEmail(email string) error {

	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}

	return nil

}
