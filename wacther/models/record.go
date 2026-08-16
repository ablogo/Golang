package models

import "time"

type Record struct {
	Id          int
	App         string
	Message     string
	Description string
	IsError     bool
	CreatedAt   time.Time
}
