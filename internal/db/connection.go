package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection() {
	// Tomamos valores del entorno para que funcione en Docker
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var err error
	// Retry simple: espera hasta que la DB esté lista
	for range 10 {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Base de datos conectada")
			return
		}
		log.Printf("Error conectando a la DB: %v, reintentando en 2s...", err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("No se pudo conectar a la base de datos después de varios intentos")
}
