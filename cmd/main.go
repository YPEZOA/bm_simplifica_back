package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/middleware"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/services/auth"
	"github.com/ypezoa/bm-simplifica-back/internal/services/company"
	"github.com/ypezoa/bm-simplifica-back/internal/services/file"
	"github.com/ypezoa/bm-simplifica-back/internal/services/user"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db.DBConnection()

	// migrate models
	err := db.DB.AutoMigrate(&models.User{}, &models.File{}, &models.Company{})
	if err != nil {
		log.Fatalf("Error migrando tablas: %v", err)
	}

	// [DEV] Seed usuarios de desarrollo
	seedDevUsers()

	// router
	r := mux.NewRouter()

	// Apply global middlewares
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// Public routes (no auth required)
	auth.AuthRoutes(r)

	// Protected routes (auth required)
	user.UserRoutes(r)
	company.CompanyRoutes(r)
	file.FileRoutes(r)

	// Get server port from environment variable, default to 8080
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Service on port :%s", port)
	http.ListenAndServe(":"+port, r)
}

func seedDevUsers() {
	env := os.Getenv("ENV")
	if env != "development" {
		return
	}

	log.Println("🔧 [DEV] Iniciando seed de usuarios de desarrollo...")

	seedUser("Administrator", "admin@simplifica.com", "Admin123!", models.RoleAdmin)
	seedUser("Cliente Demo", "cliente@simplifica.com", "Cliente123!", models.RoleUser)

	log.Println("✅ [DEV] Seed de desarrollo completado")
}

func seedUser(name, email, password string, role models.Role) {
	var existingUser models.User
	result := db.DB.Where("email = ?", email).First(&existingUser)

	if result.Error == nil {
		log.Printf("⚠️  [DEV] Usuario %s ya existe, omitiendo...", email)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ [DEV] Error hasheando contraseña para %s: %v", email, err)
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		Phone:    "+56912345678",
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Printf("❌ [DEV] Error creando usuario %s: %v", email, err)
		return
	}

	fmt.Printf("✅ [DEV] Usuario creado: %s | %s | %s\n", email, password, role)
}
