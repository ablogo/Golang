package models

import "time"

type User struct {
	Id               int `gorm:"primaryKey"`
	Name             string
	LastName         string
	Email            string
	EmailVerified    bool
	Password         string
	Issuer           string
	Online           bool `gorm:"default:NULL"`
	TwoFactorEnabled bool `gorm:"default:NULL"`
	Disabled         bool
	CreatedAt        time.Time
	UpdatedAt        *time.Time `gorm:"default:NULL"`
	PictureId        *int
	Picture          *Picture `gorm:"foreignKey:PictureId"`
	Address          *[]Address
	Roles            *[]Roles `gorm:"many2many:user_roles;"`
}
