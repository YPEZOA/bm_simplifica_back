package types

import (
	"github.com/google/uuid"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStorage interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uuid.UUID) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	DeleteUser(id uuid.UUID) error
	UpdateUserPassword(id uuid.UUID, hashedPassword string) error
}

type FileStorage interface {
	GetFiles(userID, companyID uuid.UUID) ([]models.File, error)
	UploadFile(file models.File) (models.File, error)
	DeleteFile(id uuid.UUID) (models.File, error)
	DeleteMultipleFiles(ids []uuid.UUID) ([]models.File, error)
}

type AuthStorage interface {
	SignIn(email string, password string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
}

type CompanyStorage interface {
	CreateCompany(company models.Company) (models.Company, error)
	GetCompanies() ([]models.Company, error)
	GetCompanyByID(id uuid.UUID) (models.Company, error)
	DeleteCompany(id uuid.UUID) error
}
