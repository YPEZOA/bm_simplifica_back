package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "client"
)

type User struct {
	gorm.Model

	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique"`
	Password  string    `gorm:"not null"`
	Role      Role      `gorm:"not null"`
	Phone     string
	Companies []Company `gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	Phone     string    `json:"phone"`
	Companies []Company `json:"companies"`
}
