package file

import (
	"github.com/google/uuid"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

type storage struct{}

func NewFileStorage() types.FileStorage {
	return &storage{}
}

func (s *storage) GetFiles(userID, companyID uuid.UUID) ([]models.File, error) {
	files := []models.File{}

	err := db.DB.Joins("JOIN companies ON companies.id = files.company_id").Where("companies.user_id = ? AND companies.id = ?", userID, companyID).Find(&files).Error
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *storage) UploadFile(file models.File) (models.File, error) {
	if err := db.DB.Create(&file).Error; err != nil {
		return models.File{}, err
	}
	return file, nil
}

func (s *storage) DeleteFile(id uuid.UUID) (models.File, error) {
	var file models.File
	if err := db.DB.First(&file, id).Error; err != nil {
		return models.File{}, err
	}

	if err := db.DB.Delete(&file).Error; err != nil {
		return models.File{}, err
	}
	return file, nil
}
