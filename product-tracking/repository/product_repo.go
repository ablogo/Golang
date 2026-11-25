package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"product-tracking/model"
	"time"

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
	query := ` INSERT INTO products(name, url, tag_id, initial_price, current_price, lowest_price, image_url)
	           VALUES ($1, $2, $3, $4, $5, $6, $7)
			   RETURNING id, created_at`
	err := r.DB.QueryRow(ctx, query, product.Name, product.Url, product.TagId, product.InitialPrice, product.CurrentPrice, product.LowestPrice, product.ImageUrl).Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		slog.Error("Fail to create product.", "description", err.Error())
		return err
	}
	return nil
}

func (r *ProductRepo) GetById(ctx context.Context, id int) *model.Product {
	query := `SELECT id, name, url, initial_price, current_price, lowest_price, tag_id, image_url, is_tracked, created_at FROM products WHERE id = $1`

	var product model.Product
	err := r.DB.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Url, &product.InitialPrice, &product.CurrentPrice, &product.LowestPrice, &product.TagId, &product.ImageUrl, &product.IsTracked, &product.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Product doesn't exist in db.", "description", err.Error())
		}
		slog.Error("Fail to get product.", "description", err.Error())
		return nil
	}
	return &product
}

func (r *ProductRepo) GetHistoricalLowestPricesById(ctx context.Context, id int) *[]model.Price {
	query := `SELECT p.price, DATE(p.created_at) FROM prices p WHERE p.product_id = $1 
	          GROUP BY p.price, DATE(p.created_at) 
			  ORDER BY DATE(p.created_at);`

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Not prices in db.", "description", err.Error())
		} else {
			slog.Error("Fail to get historical prices.", "description", err.Error())
		}
		return nil
	}
	defer rows.Close()

	var prices []model.Price
	for rows.Next() {
		var p model.Price
		err := rows.Scan(&p.Price, &p.CreatedAt)
		if err != nil {
			slog.Error("Fail to scan product.", "description", err.Error())
			return nil
		}
		prices = append(prices, p)
	}
	return &prices
}

func (r *ProductRepo) GetChartInfoById(ctx context.Context, id int) *model.ChartInfo {
	query := `SELECT p.price, DATE(p.created_at) FROM prices p WHERE p.product_id = $1 
	          GROUP BY p.price, DATE(p.created_at) 
			  ORDER BY DATE(p.created_at);`

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Not prices in db.", "description", err.Error())
		} else {
			slog.Error("Fail to get historical prices.", "description", err.Error())
		}
		return nil
	}
	defer rows.Close()

	var prices []float64
	var dates []string
	for rows.Next() {
		var p model.Prices
		var d model.Dates
		err := rows.Scan(&p.Price, &d.CreatedAt)
		if err != nil {
			slog.Error("Fail to scan historical prices.", "description", err.Error())
			return nil
		}
		prices = append(prices, p.Price)
		dates = append(dates, d.CreatedAt.Format(time.DateOnly))
	}

	info := model.ChartInfo{
		Prices: prices,
		Dates:  dates,
	}

	return &info
}

func (r *ProductRepo) GetByIdWithHistoricalPrices(ctx context.Context, id int) *model.Product {
	query := `SELECT JSON_BUILD_OBJECT(
	          'id', a.id,
			  'name', a.name,
			  'url', a.url,
			  'initial_price', a.initial_price,
			  'current_price', a.current_price,
			  'lowest_price', a.lowest_price,
			  'tag_id', a.tag_id,
			  'image_url', a.image_url,
			  'is_tracked', a.is_tracked,
			  'created_at', a.created_at,
			  'prices', COALESCE (JSON_AGG(b.*) FILTER (WHERE b.price IS NOT NULL), '[]'))
			  FROM products a
			  LEFT JOIN prices b ON b.product_id = a.id
			  WHERE a.id = $1
			  GROUP BY a.id;`

	var jsonData []byte
	err := r.DB.QueryRow(ctx, query, id).Scan(&jsonData)
	if err != nil {
		slog.Error("Fail to get product from db.", "description", err.Error())
		return nil
	}

	var product model.Product
	if err := json.Unmarshal(jsonData, &product); err != nil {
		slog.Error("Fail to scan product.", "description", err.Error())
		return nil
	}

	return &product
}

func (r *ProductRepo) GetAll(ctx context.Context) *[]model.Product {
	query := `SELECT id, name, url, initial_price, current_price, lowest_price, tag_id, image_url, is_tracked, created_at FROM products ORDER BY id;`
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
		err := rows.Scan(&p.ID, &p.Name, &p.Url, &p.InitialPrice, &p.CurrentPrice, &p.LowestPrice, &p.TagId, &p.ImageUrl, &p.IsTracked, &p.CreatedAt)
		if err != nil {
			slog.Error("Fail to scan product.", "description", err.Error())
			return nil
		}
		products = append(products, p)
	}

	return &products
}

func (r *ProductRepo) GetAllToTrack(ctx context.Context) *[]model.Product {
	query := `SELECT id, name, url, initial_price, current_price, lowest_price, tag_id, image_url, is_tracked, created_at FROM products WHERE is_tracked = True;`
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
		err := rows.Scan(&p.ID, &p.Name, &p.Url, &p.InitialPrice, &p.CurrentPrice, &p.LowestPrice, &p.TagId, &p.ImageUrl, &p.IsTracked, &p.CreatedAt)
		if err != nil {
			slog.Error("Fail to scan product.", "description", err.Error())
			return nil
		}
		products = append(products, p)
	}

	return &products
}

func (r *ProductRepo) GetHistoricalLowestPrices(ctx context.Context) *[]model.Price {
	query := `SELECT DISTINCT ON (DATE(p.created_at)) p.current_price, p.created_at FROM prices as p WHERE p.product_id = $1
	          AND p.current_price = (SELECT MIN(p1.price) FROM prices as p1 WHERE p1.product_id = $1 and DATE(p1.created_at) = DATE(p.created_at))
			  ORDER BY DATE(p.created_at);`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		} else {
			slog.Error("Fail to get products.", "description", err.Error())
		}
		return nil
	}
	defer rows.Close()

	var prices []model.Price
	for rows.Next() {
		var p model.Price
		err := rows.Scan(&p.Price, &p.CreatedAt)
		if err != nil {
			slog.Error("Fail to scan product.", "description", err.Error())
			return nil
		}
		prices = append(prices, p)
	}

	return &prices
}

func (r *ProductRepo) GetAllWithHistoricalPrices(ctx context.Context) *[]model.Product {
	query := `SELECT a.id, a.name, a.url, a.initial_price, a.current_price, a.lowest_price, a.tag_id, a.image_url, a.is_tracked, a.created_at,
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
		err := rows.Scan(&p.ID, &p.Name, &p.Url, &p.InitialPrice, &p.CurrentPrice, &p.LowestPrice, &p.TagId, &p.ImageUrl, &p.IsTracked, &p.CreatedAt, &pricesJSON)
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

func (r *ProductRepo) ChangeTrackedStatus(ctx context.Context, product_id int) error {
	query := `UPDATE products SET is_tracked = 
	          (SELECT CASE WHEN is_tracked IS TRUE THEN FALSE WHEN is_tracked IS FALSE THEN TRUE END AS tracked FROM products WHERE id = $1) WHERE id = $1`

	commandTag, err := r.DB.Exec(ctx, query, product_id)
	if err != nil {
		slog.Error("Fail to update product.", "description", err.Error())
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("Product not found for update")
	}
	return nil
}

func (r *ProductRepo) Update(ctx context.Context, product *model.Product) error {
	query := `UPDATE products SET name = $1, url = $2, tag_id = $3, current_price = $4 WHERE id = $5`

	commandTag, err := r.DB.Exec(ctx, query, product.Name, product.Url, product.TagId, product.CurrentPrice, product.ID)
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
