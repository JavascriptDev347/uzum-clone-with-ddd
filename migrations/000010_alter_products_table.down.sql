ALTER TABLE products
    DROP CONSTRAINT IF EXISTS products_slug_key,
    DROP CONSTRAINT IF EXISTS products_packaging_type_check,
    DROP CONSTRAINT IF EXISTS products_freshness_lifespan_check,
    DROP CONSTRAINT IF EXISTS products_rating_check;

ALTER TABLE products
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS images,
    DROP COLUMN IF EXISTS video_url,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS discount_price_amount,
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS is_available,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS stock,
    DROP COLUMN IF EXISTS sold_count,
    DROP COLUMN IF EXISTS flower_types,
    DROP COLUMN IF EXISTS color,
    DROP COLUMN IF EXISTS stem_count,
    DROP COLUMN IF EXISTS packaging_type,
    DROP COLUMN IF EXISTS freshness_lifespan,
    DROP COLUMN IF EXISTS care_instructions,
    DROP COLUMN IF EXISTS occasions,
    DROP COLUMN IF EXISTS allow_custom_card,
    DROP COLUMN IF EXISTS compatible_addons;

ALTER TABLE products
    ADD COLUMN image_url TEXT,
    ADD COLUMN image_public_id TEXT;
