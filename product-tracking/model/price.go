package model

import "time"

type Price struct {
	ID        int       `db:"id" json:"id"`
	Price     float64   `db:"price" json:"price"`
	TagText   string    `db:"tag_text" json:"tag_text"`
	ProductID int       `db:"product_id" json:"product_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
