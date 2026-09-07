package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (db *gorm.DB) {
	dsn := "host=144.217.166.230 user=ubuntu-postgres dbname=auth password=Berenice1108. port=32012 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting db")
	}
	return
}
