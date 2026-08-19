ALTER TABLE products
    RENAME COLUMN video_url TO video_url_instagram;

ALTER TABLE products
    ADD COLUMN video_url_youtube TEXT NOT NULL DEFAULT '';
