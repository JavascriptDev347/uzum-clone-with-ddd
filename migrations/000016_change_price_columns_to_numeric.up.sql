ALTER TABLE products
    ALTER COLUMN price_amount TYPE NUMERIC(15,2) USING price_amount::numeric(15,2),
    ALTER COLUMN discount_price_amount TYPE NUMERIC(15,2) USING discount_price_amount::numeric(15,2);
