ALTER TABLE categories
    ADD COLUMN name_uz  VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN name_eng VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN name_ru  VARCHAR(255) NOT NULL DEFAULT '';

-- mavjud qatorlar bo'lsa, bitta tilda saqlangan eski name barcha tillarga nusxalanadi
UPDATE categories SET name_uz = name, name_eng = name, name_ru = name;

ALTER TABLE categories
    DROP COLUMN name;

CREATE INDEX ON categories(name_uz);
CREATE INDEX ON categories(name_eng);
CREATE INDEX ON categories(name_ru);
