package service

import (
	"context"
	"log/slog"
	"regexp"
	"strconv"

	"product-tracking/model"
	"product-tracking/repository"
)

type ProductService struct {
	Repo repository.ProductRepository
	Ctx  context.Context
}

func (s *ProductService) Create(new_product model.Product) (product *model.Product) {
	db_result := s.Repo.Create(s.Ctx, &new_product)

	if db_result != nil {
		return nil
	}
	return &new_product
}

func (s *ProductService) RegisterPrice(product_id int, website_text string) (float64, error) {
	re := regexp.MustCompile(`\d+(?:\.\d+)?%`)
	website_price := re.ReplaceAllString(website_text, "")

	reg := regexp.MustCompile(`[^\d.-]`)
	cleaned_text := reg.ReplaceAllString(website_price, "")

	cleaned_price, err := strconv.ParseFloat(cleaned_text, 32)
	if err != nil {
		slog.Error("Scraped price is not a valid number.", "product_id", product_id)
	}

	price := model.Price{
		ProductID: product_id,
		Price:     cleaned_price,
		TagText:   website_text,
	}

	db_result := s.Repo.RegisterPrice(s.Ctx, &price)
	return cleaned_price, db_result
}

func (s *ProductService) GetById(id int, history bool) (product *model.Product) {
	if history {
		product = s.Repo.GetByIdWithHistoricalPrices(s.Ctx, id)
	} else {
		product = s.Repo.GetById(s.Ctx, id)
	}
	return
}

func (s *ProductService) GetProducts(history bool) (products *[]model.Product) {
	if history {
		products = s.Repo.GetAllWithHistoricalPrices(s.Ctx)
	} else {
		products = s.Repo.GetAll(s.Ctx)
	}
	return
}

func (s *ProductService) IsLowestPrice(id int, price float64) bool {
	result := s.Repo.CallProcedure(s.Ctx, id, price)
	return result
}

func (s *ProductService) DeleteProduct(id int) bool {
	result := s.Repo.Delete(s.Ctx, id)

	if result != nil {
		return true
	} else {
		return false
	}
}

func (s *ProductService) UpdateProduct(updated_product model.Product) bool {
	result := s.Repo.Update(s.Ctx, &updated_product)

	if result != nil {
		return true
	} else {
		return false
	}
}

func (s *ProductService) UpdateProductPrice(Id int, price int) bool {
	result := s.Repo.UpdatePrice(s.Ctx, Id, price)

	if result != nil {
		return true
	} else {
		return false
	}
}
