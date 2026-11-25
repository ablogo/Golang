package model

import "time"

type Price struct {
	ID        int       `db:"id" json:"id,omitempty"`
	Price     float64   `db:"price" json:"price"`
	TagText   string    `db:"tag_text" json:"tag_text,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	ProductID int       `db:"product_id" json:"-"`
}
