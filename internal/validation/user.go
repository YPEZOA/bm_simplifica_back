package validation

import (
	"errors"
	"regexp"
	"strings"
)

type UserValidator struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

func (uv *UserValidator) Validate() error {
	if err := uv.validateName(); err != nil {
		return err
	}
	if err := uv.validateEmail(); err != nil {
		return err
	}
	if err := uv.validatePassword(); err != nil {
		return err
	}
	if err := uv.validatePhone(); err != nil {
		return err
	}
	return nil
}

func (uv *UserValidator) validateName() error {
	name := strings.TrimSpace(uv.Name)
	if len(name) < 2 {
		return errors.New("El nombre debe tener al menos 2 caracteres")
	}
	if len(name) > 100 {
		return errors.New("El nombre debe tener menos de 100 caracteres")
	}
	return nil
}

func (uv *UserValidator) validateEmail() error {
	email := strings.TrimSpace(uv.Email)
	if email == "" {
		return errors.New("El email es requerido")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Formato de email incorrecto")
	}
	return nil
}

func (uv *UserValidator) validatePassword() error {
	password := uv.Password
	if len(password) < 8 {
		return errors.New("La contraseña debe tener al menos 8 caracteres")
	}
	if len(password) > 128 {
		return errors.New("La contraseña debe tener menos de 128 caracteres")
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	if !hasUpper {
		return errors.New("La contraseña debe contener al menos una letra mayúscula")
	}
	if !hasLower {
		return errors.New("La contraseña debe contener al menos una letra minúscula")
	}
	if !hasNumber {
		return errors.New("La contraseña debe contener al menos un número")
	}
	if !hasSpecial {
		return errors.New("La contraseña debe contener al menos un carácter especial")
	}

	return nil
}

func (uv *UserValidator) validatePhone() error {
	phone := strings.TrimSpace(uv.Phone)
	if phone == "" {
		return nil
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("Formato de número de teléfono inválido")
	}
	return nil
}

type LoginValidator struct {
	Email    string
	Password string
}

func (lv *LoginValidator) Validate() error {
	if err := lv.validateEmail(); err != nil {
		return err
	}
	if err := lv.validatePassword(); err != nil {
		return err
	}
	return nil
}

func (lv *LoginValidator) validateEmail() error {
	email := strings.TrimSpace(lv.Email)
	if email == "" {
		return errors.New("El email es requerido")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Formato de email incorrecto")
	}
	return nil
}

func (lv *LoginValidator) validatePassword() error {
	password := lv.Password
	if password == "" {
		return errors.New("La contraseña es requerida")
	}
	return nil
}

func (uv *UserValidator) ValidatePassword() error {
	return uv.validatePassword()
}
