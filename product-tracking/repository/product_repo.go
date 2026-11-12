package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"product-tracking/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &ProductRepo{DB: db}
}

func (r *ProductRepo) Create(ctx context.Context, product *model.Product) error {
	query := ` INSERT INTO products(name, url)
	           VALUES ($1, $2)
			   RETURNING id, created_at`
	err := r.DB.QueryRow(ctx, query, product.Name, product.Url).Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		slog.Error("Fail to create product.", "description", err.Error())
		return err
	}
	return nil
}

func (r *ProductRepo) GetById(ctx context.Context, id int) *model.Product {
	query := `SELECT id, name, url, price, tag_id, created_at FROM products WHERE id = $1`

	var product model.Product
	err := r.DB.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Url, &product.Price, &product.TagId, &product.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		}
		slog.Error("Fail to get product.", "description", err.Error())
		return nil
	}
	return &product
}

func (r *ProductRepo) GetByIdWithHistoricalPrices(ctx context.Context, id int) *model.Product {
	query := `SELECT JSON_BUILD_OBJECT(
	          'id', a.id,
			  'name', a.name,
			  'url', a.url,
			  'price', a.price,
			  'tag_id', a.tag_id,
			  'created_at', a.created_at,
			  'prices', COALESCE (JSON_AGG(b.*) FILTER (WHERE b.price IS NOT NULL), '[]'))
			  FROM products a
			  LEFT JOIN prices b ON b.product_id = a.id
			  WHERE a.id = $1;`

	var jsonData []byte
	err := r.DB.QueryRow(ctx, query, id).Scan(&jsonData)
	if err != nil {
		return nil
	}

	var product model.Product
	if err := json.Unmarshal(jsonData, &product); err != nil {
		return nil
	}

	return &product
}

func (r *ProductRepo) GetAll(ctx context.Context) *[]model.Product {
	query := `SELECT id, name, url, price, tag_id, created_at FROM products ORDER BY id;`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		} else {
			slog.Error("Fail to get products.", "description", err.Error())
		}
		return nil
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Url, &p.Price, &p.TagId, &p.CreatedAt)
		if err != nil {
			slog.Error("Fail to scan product.", "description", err.Error())
			return nil
		}
		products = append(products, p)
	}

	return &products
}

func (r *ProductRepo) GetAllWithHistoricalPrices(ctx context.Context) *[]model.Product {
	query := `SELECT a.id, a.name, a.url, a.price, a.tag_id, a.created_at,
	          COALESCE (JSON_AGG(JSON_BUILD_OBJECT(
			  'id', b.id, 
			  'price', b.price,
			  'tag_text', b.tag_text,
			  'created_at', b.created_at)) FILTER (WHERE b.price IS NOT NULL), '[]') 
			  AS prices
			  FROM products a
			  LEFT JOIN prices b ON b.product_id = a.id
			  GROUP BY a.id
			  ORDER BY a.id;`

	var products []model.Product
	var pricesJSON []byte

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		} else {
			slog.Error("Fail to get products.", "description", err.Error())
		}
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Url, &p.Price, &p.TagId, &p.CreatedAt, &pricesJSON)
		if err != nil {
			slog.Error("Fail to scan product.", "description", err.Error())
			return nil
		}
		if err := json.Unmarshal(pricesJSON, &p.Prices); err != nil {
			slog.Error("Fail to scan prices.", "description", err.Error())
		}
		products = append(products, p)
	}

	return &products
}

func (r *ProductRepo) RegisterPrice(ctx context.Context, record *model.Price) error {
	query := ` INSERT INTO prices(price, tag_text, product_id)
	           VALUES ($1, $2, $3)
			   RETURNING id, created_at`
	err := r.DB.QueryRow(ctx, query, record.Price, record.TagText, record.ProductID).Scan(&record.ID, &record.CreatedAt)
	if err != nil {
		slog.Error("Fail to save price.", "description", err.Error())
		return err
	}
	return nil
}

func (r *ProductRepo) Update(ctx context.Context, product *model.Product) error {
	query := `UPDATE product SET name = $1, url = $2 WHERE id = $3`

	commandTag, err := r.DB.Exec(ctx, query, product.Name, product.Url, product.ID)
	if err != nil {
		slog.Error("Fail to update product.", "description", err.Error())
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("Product not found for update")
	}
	return nil
}

func (r *ProductRepo) UpdatePrice(ctx context.Context, id int, price int) error {
	query := `UPDATE product SET price = $1 WHERE id = $2`

	commandTag, err := r.DB.Exec(ctx, query, price, id)
	if err != nil {
		slog.Error("Fail to update product price.", "description", err.Error())
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("Product not found for update")
	}
	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Fail to delete product.", "description", err.Error())
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("Product not found for deletion")
	}
	return nil
}

func (r *ProductRepo) CallProcedure(ctx context.Context, product_id int, price float64) bool {
	query := `CALL is_lowest_price($1, $2, NULL);`

	var isLowestPrice bool
	err := r.DB.QueryRow(ctx, query, product_id, price).Scan(&isLowestPrice)
	if err != nil {
		slog.Error("Fail to get product.", "description", err.Error())
	}
	return isLowestPrice
}
