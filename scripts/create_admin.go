package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file from current directory
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Connect to database
	db.DBConnection()

	// Create admin user
	adminUser := models.User{
		ID:    uuid.New(),
		Name:  "Administrator",
		Email: "admin@simplifica.com",
		Role:  models.RoleAdmin,
		Phone: "+56998765432",
	}

	// Hash password
	password := "Admin123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error al encriptar contraseña: %v", err)
	}
	adminUser.Password = string(hashedPassword)

	// Create user in database
	result := db.DB.Create(&adminUser)
	if result.Error != nil {
		log.Fatalf("Error al crear usuario admin: %v", result.Error)
	}

	fmt.Printf("✅ Usuario admin creado exitosamente:\n")
	fmt.Printf("📧 Email: %s\n", adminUser.Email)
	fmt.Printf("🔑 Contraseña: %s\n", password)
	fmt.Printf("👤 Rol: %s\n", adminUser.Role)
	fmt.Printf("🆔 ID: %s\n", adminUser.ID)

}
