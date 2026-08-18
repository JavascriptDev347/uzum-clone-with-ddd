CREATE TABLE events (
    id UUID PRIMARY KEY,
    eyebrow VARCHAR(255) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL DEFAULT '',
    cta VARCHAR(255) NOT NULL DEFAULT '',
    image_url TEXT NOT NULL,
    image_public_id TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    is_root BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX ON events(deleted_at);
CREATE INDEX ON events(category_id);
CREATE INDEX ON events(is_root);
