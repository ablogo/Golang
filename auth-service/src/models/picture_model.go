package models

import "time"

type Picture struct {
	Id          int     `gorm:"primaryKey"`
	Picture     *[]byte `gorm:"default:NULL"`
	FileName    *string `gorm:"default:NULL"`
	ContentType *string `gorm:"default:NULL"`
	PictureUrl  *string `gorm:"default:NULL"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time `gorm:"default:NULL"`
}
