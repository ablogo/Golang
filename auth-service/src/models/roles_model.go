package models

type Roles struct {
	Id   int `gorm:"primaryKey"`
	Name string
}
