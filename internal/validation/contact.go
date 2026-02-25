package validation

import (
	"errors"
	"regexp"
	"strings"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func (cr *ContactRequest) Validate() error {
	if err := cr.validateName(); err != nil {
		return err
	}
	if err := cr.validateEmail(); err != nil {
		return err
	}
	if err := cr.validatePhone(); err != nil {
		return err
	}
	if err := cr.validateMessage(); err != nil {
		return err
	}
	return nil
}

func (cr *ContactRequest) validateName() error {
	name := strings.TrimSpace(cr.Name)
	if len(name) < 2 {
		return errors.New("El nombre debe tener al menos 2 caracteres")
	}
	if len(name) > 100 {
		return errors.New("El nombre debe tener menos de 100 caracteres")
	}
	return nil
}

func (cr *ContactRequest) validateEmail() error {
	email := strings.TrimSpace(cr.Email)
	if email == "" {
		return errors.New("El email es requerido")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Formato de email incorrecto")
	}
	return nil
}

func (cr *ContactRequest) validatePhone() error {
	phone := strings.TrimSpace(cr.Phone)
	if phone == "" {
		return errors.New("El número de teléfono es requerido")
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("Formato de número de teléfono inválido")
	}
	return nil
}

func (cr *ContactRequest) validateMessage() error {
	message := strings.TrimSpace(cr.Message)
	if len(message) < 10 {
		return errors.New("El mensaje debe tener al menos 10 caracteres")
	}
	if len(message) > 1000 {
		return errors.New("El mensaje debe tener menos de 1000 caracteres")
	}
	return nil
}
