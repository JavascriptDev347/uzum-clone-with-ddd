ALTER TABLE events
    ADD COLUMN eyebrow  VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN title    VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN subtitle VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN cta      VARCHAR(255) NOT NULL DEFAULT '';

UPDATE events SET
    eyebrow = eyebrow_uz, title = title_uz, subtitle = subtitle_uz, cta = cta_uz;

ALTER TABLE events
    DROP COLUMN eyebrow_uz,
    DROP COLUMN eyebrow_eng,
    DROP COLUMN eyebrow_ru,
    DROP COLUMN title_uz,
    DROP COLUMN title_eng,
    DROP COLUMN title_ru,
    DROP COLUMN subtitle_uz,
    DROP COLUMN subtitle_eng,
    DROP COLUMN subtitle_ru,
    DROP COLUMN cta_uz,
    DROP COLUMN cta_eng,
    DROP COLUMN cta_ru;
