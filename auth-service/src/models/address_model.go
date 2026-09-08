package models

import "time"

type Address struct {
	Id        int32 `gorm:"primaryKey"`
	Country   string
	State     string
	Colony    string
	Street    string
	Number    string
	Disabled  bool
	CreatedAt time.Time
	UpdatedAt *time.Time `gorm:"default:NULL"`
	UserId    int
}
