ALTER TABLE categories
    ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';

UPDATE categories SET name = name_uz;

DROP INDEX IF EXISTS categories_name_uz_idx;
DROP INDEX IF EXISTS categories_name_eng_idx;
DROP INDEX IF EXISTS categories_name_ru_idx;

ALTER TABLE categories
    DROP COLUMN name_uz,
    DROP COLUMN name_eng,
    DROP COLUMN name_ru;
