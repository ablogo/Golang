package repository

import (
	"context"
	"product-tracking/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetById(ctx context.Context, id int) *model.Product
	GetByIdWithHistoricalPrices(ctx context.Context, id int) *model.Product
	GetHistoricalLowestPricesById(ctx context.Context, id int) *[]model.Price
	GetAll(ctx context.Context) *[]model.Product
	GetAllWithHistoricalPrices(ctx context.Context) *[]model.Product
	GetHistoricalLowestPrices(ctx context.Context) *[]model.Price
	Update(ctx context.Context, product *model.Product) error
	UpdatePrice(ctx context.Context, id int, price int) error
	Delete(ctx context.Context, id int) error
	RegisterPrice(ctx context.Context, price *model.Price) error
	CallProcedure(ctx context.Context, product_id int, price float64) bool
}
