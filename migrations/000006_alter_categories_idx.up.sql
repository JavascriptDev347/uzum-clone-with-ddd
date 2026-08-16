CREATE UNIQUE INDEX idx_categories_parent_name ON categories (parent_id, name) WHERE deleted_at IS NULL;
