ALTER TABLE products
    ALTER COLUMN price_amount TYPE BIGINT USING price_amount::bigint,
    ALTER COLUMN discount_price_amount TYPE BIGINT USING discount_price_amount::bigint;
