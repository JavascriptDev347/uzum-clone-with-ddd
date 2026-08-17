ALTER TABLE categories
    ADD COLUMN image_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN image_public_id TEXT NOT NULL DEFAULT '';

ALTER TABLE categories ALTER COLUMN image_url DROP DEFAULT;
ALTER TABLE categories ALTER COLUMN image_public_id DROP DEFAULT;
