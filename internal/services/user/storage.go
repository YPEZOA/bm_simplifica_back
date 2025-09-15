package user

import (
	"github.com/google/uuid"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

type storage struct{}

func NewUserStorage() types.UserStorage {
	return &storage{}
}

func (s *storage) GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *storage) GetUserByID(id uuid.UUID) (models.User, error) {
	user := models.User{}
	if err := db.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (s *storage) CreateUser(user models.User) (models.User, error) {
	if err := db.DB.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
