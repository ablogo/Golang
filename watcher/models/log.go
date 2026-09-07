package models

import "time"

type Log struct {
	Id         int
	Level      string
	Message    string
	Source     string
	Attributes string
	CreatedAt  time.Time
}
