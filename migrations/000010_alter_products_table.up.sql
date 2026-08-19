ALTER TABLE products
    DROP COLUMN IF EXISTS image_url,
    DROP COLUMN IF EXISTS image_public_id;

ALTER TABLE products
    ADD COLUMN description        TEXT NOT NULL DEFAULT '',
    ADD COLUMN images              JSONB NOT NULL DEFAULT '[]',
    ADD COLUMN video_url           TEXT NOT NULL DEFAULT '',
    ADD COLUMN updated_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN discount_price_amount BIGINT,
    ADD COLUMN slug                VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN is_available        BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN rating              NUMERIC(2,1) NOT NULL DEFAULT 1,
    ADD COLUMN stock               INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN sold_count          INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN flower_types        TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN color               VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN stem_count          INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN packaging_type      VARCHAR(20) NOT NULL DEFAULT 'box',
    ADD COLUMN freshness_lifespan  INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN care_instructions   TEXT,
    ADD COLUMN occasions           TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN allow_custom_card   BOOLEAN,
    ADD COLUMN compatible_addons   TEXT[] NOT NULL DEFAULT '{}';

-- mavjud qatorlar bo'lsa, unique constraint qo'shishdan oldin slug'ni to'ldiramiz
UPDATE products SET slug = id::text WHERE slug = '' OR slug IS NULL;

ALTER TABLE products
    ADD CONSTRAINT products_slug_key UNIQUE (slug);

ALTER TABLE products
    ADD CONSTRAINT products_packaging_type_check CHECK (packaging_type IN ('bucket', 'box', 'vase'));

ALTER TABLE products
    ADD CONSTRAINT products_freshness_lifespan_check CHECK (freshness_lifespan BETWEEN 1 AND 7);

ALTER TABLE products
    ADD CONSTRAINT products_rating_check CHECK (rating BETWEEN 1 AND 5);

CREATE INDEX ON products(slug);
CREATE INDEX ON products(is_available);
