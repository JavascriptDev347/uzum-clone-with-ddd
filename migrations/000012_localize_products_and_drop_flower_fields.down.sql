ALTER TABLE products
    ADD COLUMN name                VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN description         TEXT NOT NULL DEFAULT '',
    ADD COLUMN video_url_youtube   TEXT NOT NULL DEFAULT '',
    ADD COLUMN video_url_instagram TEXT NOT NULL DEFAULT '',
    ADD COLUMN flower_types        TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN color               VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN stem_count          INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN packaging_type      VARCHAR(20) NOT NULL DEFAULT 'box',
    ADD COLUMN freshness_lifespan  INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN care_instructions   TEXT,
    ADD COLUMN allow_custom_card   BOOLEAN,
    ADD COLUMN compatible_addons   TEXT[] NOT NULL DEFAULT '{}';

UPDATE products SET name = name_uz, description = description_uz;

ALTER TABLE products
    ADD CONSTRAINT products_packaging_type_check CHECK (packaging_type IN ('bucket', 'box', 'vase')),
    ADD CONSTRAINT products_freshness_lifespan_check CHECK (freshness_lifespan BETWEEN 1 AND 7);

DROP INDEX IF EXISTS products_name_uz_idx;
DROP INDEX IF EXISTS products_name_eng_idx;
DROP INDEX IF EXISTS products_name_ru_idx;

ALTER TABLE products
    DROP COLUMN name_uz,
    DROP COLUMN name_eng,
    DROP COLUMN name_ru,
    DROP COLUMN description_uz,
    DROP COLUMN description_eng,
    DROP COLUMN description_ru,
    DROP COLUMN tag_uz,
    DROP COLUMN tag_eng,
    DROP COLUMN tag_ru;
