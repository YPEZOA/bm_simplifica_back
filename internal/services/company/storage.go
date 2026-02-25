package company

import (
	"github.com/google/uuid"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
	"gorm.io/gorm"
)

type storage struct{}

func NewCompanyStorage() types.CompanyStorage {
	return &storage{}
}

func (*storage) GetCompanies() ([]models.Company, error) {
	companies := []models.Company{}
	if err := db.DB.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (*storage) CreateCompany(company models.Company) (models.Company, error) {
	// Verificar que el usuario existe antes de crear la compañía
	var userExists models.User
	if err := db.DB.Select("id").Where("id = ?", company.UserID).First(&userExists).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return company, gorm.ErrRecordNotFound
		}
		return company, err
	}

	// Crear la compañía
	if err := db.DB.Create(&company).Error; err != nil {
		return company, err
	}
	return company, nil
}

func (*storage) GetCompanyByID(id uuid.UUID) (models.Company, error) {
	company := models.Company{}
	if err := db.DB.First(&company, id).Error; err != nil {
		return company, err
	}
	return company, nil
}

func (*storage) DeleteCompany(id uuid.UUID) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Primero eliminar archivos de la compañía
		if err := tx.Where("company_id = ?", id).Delete(&models.File{}).Error; err != nil {
			return err
		}

		// Luego eliminar la compañía
		if err := tx.Delete(&models.Company{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}
