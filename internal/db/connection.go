package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DSN = "host=localhost user=ypezoa password=hh9m3m34 dbname=bm_simplifica port=5432 sslmode=disable"
	DB  *gorm.DB
)

func DBConnection() {
	var err error

	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Base de datos conectada")
}
