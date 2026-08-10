package postgres

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresCategoryRepository struct {
	db *sqlx.DB
}

func NewPostgresCategoryRepository(db *sqlx.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (id, name, parent_id, created_at, deleted_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, category.id, category.name, category.parentID, category.createdAt, category.deletedAt)
	return err
}
