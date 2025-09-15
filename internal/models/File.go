package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"not null"`
	Path      string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	CompanyID uuid.UUID
}
