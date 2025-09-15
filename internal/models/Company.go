package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model

	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name   string    `gorm:"not null"`
	Rut    string    `gorm:"unique, not null"`
	Files  []File    `gorm:"foreignKey:CompanyID"`
	UserID uuid.UUID
}
