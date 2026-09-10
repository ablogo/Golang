package database

import (
	"fmt"
	"log"
	"src/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBHandler struct {
	Settings *utils.Settings
}

var count = 0

func (h *DBHandler) InitDB() (db *gorm.DB) {
	count += 1
	fmt.Printf("calling db constructur: %d", count)
	// strings in golang are immutable, string builder have better performance for concatenation
	// but it is only few concatenations is ok the code is more clear, put everything in one line will be the best solution
	dsn := "host=" + h.Settings.DB_HOST
	dsn += " user=" + h.Settings.DB_USER
	dsn += " dbname=" + h.Settings.DB_NAME
	dsn += " password=" + h.Settings.DB_PASSWORD
	dsn += " port=" + h.Settings.DB_PORT + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error obtaining configuration values")
	}
	return db.Session(&gorm.Session{})
}
