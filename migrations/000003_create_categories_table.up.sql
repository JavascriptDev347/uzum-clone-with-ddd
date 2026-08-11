CREATE TABLE categories (
    id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	parent_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP
);

create index on categories(parent_id);
create index on categories(deleted_at);
