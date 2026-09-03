ALTER TABLE events
    ADD COLUMN eyebrow_uz    VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN eyebrow_eng   VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN eyebrow_ru    VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN title_uz      VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN title_eng     VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN title_ru      VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN subtitle_uz   VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN subtitle_eng  VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN subtitle_ru   VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN cta_uz        VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN cta_eng       VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN cta_ru        VARCHAR(255) NOT NULL DEFAULT '';

-- mavjud qatorlar bo'lsa, bitta tilda saqlangan eski qiymatlar barcha tillarga nusxalanadi
UPDATE events SET
    eyebrow_uz = eyebrow, eyebrow_eng = eyebrow, eyebrow_ru = eyebrow,
    title_uz = title, title_eng = title, title_ru = title,
    subtitle_uz = subtitle, subtitle_eng = subtitle, subtitle_ru = subtitle,
    cta_uz = cta, cta_eng = cta, cta_ru = cta;

ALTER TABLE events
    DROP COLUMN eyebrow,
    DROP COLUMN title,
    DROP COLUMN subtitle,
    DROP COLUMN cta;
