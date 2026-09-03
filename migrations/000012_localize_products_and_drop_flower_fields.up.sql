ALTER TABLE products
    ADD COLUMN name_uz          VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN name_eng         VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN name_ru          VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN description_uz   TEXT NOT NULL DEFAULT '',
    ADD COLUMN description_eng  TEXT NOT NULL DEFAULT '',
    ADD COLUMN description_ru   TEXT NOT NULL DEFAULT '',
    ADD COLUMN tag_uz           VARCHAR(255),
    ADD COLUMN tag_eng          VARCHAR(255),
    ADD COLUMN tag_ru           VARCHAR(255);

-- mavjud qatorlar bo'lsa, bitta tilda saqlangan eski name/description barcha tillarga nusxalanadi
UPDATE products SET
    name_uz = name, name_eng = name, name_ru = name,
    description_uz = description, description_eng = description, description_ru = description;

ALTER TABLE products
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS video_url_youtube,
    DROP COLUMN IF EXISTS video_url_instagram,
    DROP CONSTRAINT IF EXISTS products_packaging_type_check,
    DROP CONSTRAINT IF EXISTS products_freshness_lifespan_check,
    DROP COLUMN IF EXISTS flower_types,
    DROP COLUMN IF EXISTS color,
    DROP COLUMN IF EXISTS stem_count,
    DROP COLUMN IF EXISTS packaging_type,
    DROP COLUMN IF EXISTS freshness_lifespan,
    DROP COLUMN IF EXISTS care_instructions,
    DROP COLUMN IF EXISTS allow_custom_card,
    DROP COLUMN IF EXISTS compatible_addons;

CREATE INDEX ON products(name_uz);
CREATE INDEX ON products(name_eng);
CREATE INDEX ON products(name_ru);
