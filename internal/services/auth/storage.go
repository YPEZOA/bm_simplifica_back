package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

type storage struct{}

func NewAuthStorage() types.AuthStorage {
	return &storage{}
}

func (s *storage) SignIn(email, password string) (models.User, error) {
	var user models.User

	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, fmt.Errorf("credenciales incorrectas")
	}

	return user, nil
}

func (s *storage) CreateUser(user models.User) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashedPassword)
	if err := db.DB.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
