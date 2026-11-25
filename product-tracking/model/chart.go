package model

import "time"

type Prices struct {
	Price float64 `db:"price" json:"price"`
}

type Dates struct {
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ChartInfo struct {
	Prices []float64 `json:"prices"`
	Dates  []string  `json:"dates"`
}
