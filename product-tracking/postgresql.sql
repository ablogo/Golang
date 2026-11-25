CREATE DATABASE products;

\c products

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(240) NOT NULL,
    url VARCHAR(240) UNIQUE NOT NULL,
    tag_id VARCHAR(680) NOT NULL,
    image_url VARCHAR(680) NULL DEFAULT '',
    initial_price DECIMAL(12,2) NOT NULL,
    current_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    lowest_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    is_tracked BOOLEAN NOT NULL DEFAULT 1,
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

-- store procedure to get lowest price of product
CREATE OR REPLACE PROCEDURE is_lowest_price(prod_id INT, new_price FLOAT, OUT lowest BOOLEAN)
LANGUAGE plpgsql
AS $$
DECLARE
BEGIN
    IF (SELECT EXISTS (SELECT p.price FROM prices AS p WHERE p.product_id = prod_id)) = FALSE THEN
	    SELECT NOT EXISTS (
            SELECT p.price FROM prices AS p WHERE p.product_id = prod_id AND p.price > 0 AND p.price <= new_price) INTO lowest;
    ELSE
        lowest := FALSE;

    END IF;
 
END;
$$;

-- function
CREATE OR REPLACE FUNCTION verify_lowest_price()
RETURNS TRIGGER 
LANGUAGE plpgsql
AS $$
BEGIN
    IF (SELECT lowest_price FROM products WHERE id = NEW.product_id) = 0 THEN
        UPDATE products SET lowest_price = NEW.price
		WHERE id = NEW.product_id;
    ELSIF NEW.price > 0 AND NEW.price < (SELECT MIN(price) FROM prices WHERE product_id = NEW.product_id AND price > 0) THEN
	    UPDATE products SET lowest_price = NEW.price
		WHERE id = NEW.product_id;
	END IF;

    UPDATE products SET current_price = NEW.price WHERE id = NEW.product_id;
    
    RETURN NEW;
END;
$$;

-- trigger on prices table
CREATE TRIGGER verify_lowest_price_trigger
BEFORE INSERT ON prices
FOR EACH ROW
EXECUTE FUNCTION verify_lowest_price();