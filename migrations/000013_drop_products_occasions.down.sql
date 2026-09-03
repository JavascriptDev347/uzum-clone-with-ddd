ALTER TABLE products
    ADD COLUMN occasions TEXT[] NOT NULL DEFAULT '{}';
