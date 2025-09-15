package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/db"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/services/auth"
	"github.com/ypezoa/bm-simplifica-back/internal/services/user"
)

func main() {
	db.DBConnection()

	// migrate models
	err := db.DB.AutoMigrate(&models.User{}, &models.File{}, &models.Company{})
	if err != nil {
		log.Fatalf("Error migrando tablas: %v", err)
	}

	// router
	r := mux.NewRouter()
	user.UserRoutes(r)
	auth.AuthRoutes(r)

	log.Println("Service on port :8080")
	http.ListenAndServe(":8080", r)
}
