package company

import (
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
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
	if err := db.DB.Create(&company).Error; err != nil {
		return company, err
	}
	return company, nil
}
