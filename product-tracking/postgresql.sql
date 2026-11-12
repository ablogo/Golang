CREATE DATABASE products;

\c products

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(240) NOT NULL,
    url VARCHAR(240) UNIQUE NOT NULL,
    tag_id VARCHAR(280) NOT NULL,
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS prices (
    id BIGSERIAL PRIMARY KEY,
    price DECIMAL(12,2) NOT NULL,
    tag_text VARCHAR(240) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    product_id INTEGER NOT NULL,

    CONSTRAINT fk_products
        FOREIGN KEY (product_id) 
        REFERENCES products(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS logs (
    id BIGSERIAL PRIMARY KEY,
    level VARCHAR(10) NOT NULL,
    component VARCHAR(240),
    message TEXT NOT NULL,
    attributes JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_logs_timestamp ON logs(timestamp DESC);
CREATE INDEX idx_logs_level ON logs(level);

CREATE OR REPLACE PROCEDURE is_lowest_price(prod_id INT, new_price FLOAT, OUT lowest BOOLEAN)
LANGUAGE plpgsql
AS $$
DECLARE
    -- Optional local variables
BEGIN
    -- SQL and PL/pgSQL procedural statements
	SELECT NOT EXISTS (
		SELECT p.price FROM prices AS p WHERE p.product_id = prod_id and p.price <= new_price
	) INTO lowest;

    COMMIT; 
END;
$$;