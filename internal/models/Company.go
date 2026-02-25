package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model

	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name   string    `gorm:"not null" json:"name"`
	Rut    string    `gorm:"unique, not null" json:"rut"`
	Files  []File    `gorm:"foreignKey:CompanyID" json:"files"`
	UserID uuid.UUID `gorm:"type:uuid;not null"references:ID" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type CompaniesResponse struct {
	Name  string `json:"name"`
	Rut   string `json:"rut"`
	Files []File `json:"files"`
}
