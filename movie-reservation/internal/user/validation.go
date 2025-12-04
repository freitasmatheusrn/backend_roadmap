package user

import (
	"fmt"
	"net/mail"
)

func EmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("inválido %s", err)
	}
	return nil
}

func PasswordValid(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("senha deve conter ao menos 8 digítos")
	}
	return nil
}

func NameValid(name string) error {
	if len(name) <= 3 {
		return fmt.Errorf("nome deve conter ao menos 3 letras")
	}
	return nil
}
