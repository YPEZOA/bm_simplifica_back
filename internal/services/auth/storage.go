package auth

import (
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

	return user, nil
}
