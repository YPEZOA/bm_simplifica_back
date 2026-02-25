package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/middleware"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/services/auth"
	"github.com/ypezoa/bm-simplifica-back/internal/services/company"
	"github.com/ypezoa/bm-simplifica-back/internal/services/file"
	"github.com/ypezoa/bm-simplifica-back/internal/services/user"
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
