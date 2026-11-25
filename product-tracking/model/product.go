package model

import "time"

type Product struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Url          string    `db:"url" json:"url"`
	ImageUrl     string    `db:"image_url" json:"image_url"`
	TagId        string    `db:"tag_id" json:"tag_id"`
	InitialPrice float32   `db:"initial_price" json:"initial_price"`
	CurrentPrice float32   `db:"current_price" json:"current_price"`
	LowestPrice  float32   `db:"lowest_price" json:"lowest_price"`
	IsTracked    bool      `db:"is_tracked" json:"is_tracked"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	Prices       []Price   `json:"prices,omitempty"`
}
