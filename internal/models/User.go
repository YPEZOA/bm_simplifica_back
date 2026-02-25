package models

import (
	"time"

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

	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `gorm:"not null" json:"password"`
	Role      Role       `gorm:"not null" json:"role"`
	Phone     string     `json:"phone"`
	Companies []Company  `gorm:"foreignKey:UserID" json:"companies"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	Phone     string    `json:"phone"`
	Companies []Company `json:"companies"`
}
