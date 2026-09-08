package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Sqlyte struct {
	Db *sql.DB
}

func (s *Sqlyte) Initialize() error {
	query := `
	CREATE TABLE records (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		app TEXT NOT NULL,
		message TEXT NOT NULL,
		description TEXT NOT NULL,
		is_error INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := s.Db.Exec(query)
	if err != nil && !strings.Contains(err.Error(), "exists") {
		fmt.Println(time.Now().UTC(), "Error init db: ", err)
	}

	return nil
}

func (s *Sqlyte) Connect(file string) error {

	dsn := file + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(7000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("Failed to open db: %w", err)
	}
	db.SetMaxOpenConns(2)
	db.SetConnMaxLifetime(time.Hour * 24)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Database is not ready: %w", err)
	}

	s.Db = db
	return nil
}

func (s *Sqlyte) Reconnect() error {
	if err := s.Db.Ping(); err != nil {

	}

	return nil
}

func (s *Sqlyte) Close() {
	if err := s.Db.Close(); err != nil {
		fmt.Println(time.Now().UTC(), "Error closing db: ", err)
	}
	fmt.Println(time.Now().UTC(), "Database successfully closed.")
}
