package model

import "time"

type Product struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Url       string    `db:"url"`
	TagId     string    `db:"tagId"`
	Price     float32   `db:"price"`
	CreatedAt time.Time `db:"createdAt"`
	Prices    []Price   `json:"prices"`
}
