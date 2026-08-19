ALTER TABLE products
    DROP COLUMN IF EXISTS video_url_youtube;

ALTER TABLE products
    RENAME COLUMN video_url_instagram TO video_url;
