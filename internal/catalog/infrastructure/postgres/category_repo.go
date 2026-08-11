package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

	_, err := r.db.ExecContext(ctx, query, category.ID(), category.Name(), category.ParentID(), category.CreatedAt(), category.DeletedAt())
	return err
}

func (r *PostgresCategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `SELECT * FROM categories WHERE id = $1 AND deleted_at IS NULL`

	var (
		name      string
		parentID  *string
		createdAt time.Time
		deletedAt *time.Time
	)

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&name, &parentID, &createdAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return domain.NewCategoryFromRepository(id, name, parentID, createdAt, deletedAt), nil
}

func (r *PostgresCategoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, deleted_at FROM categories WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*domain.Category{}
	for rows.Next() {
		var (
			id        string
			name      string
			parentID  *string
			createdAt time.Time
			deletedAt sql.NullTime
		)
		err := rows.Scan(&id, &name, &parentID, &createdAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		var deletedAtPtr *time.Time
		if deletedAt.Valid {
			deletedAtPtr = &deletedAt.Time
		}
		categories = append(categories, domain.NewCategoryFromRepository(id, name, parentID, createdAt, deletedAtPtr))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) FindChildren(ctx context.Context, parentID string) ([]*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, deleted_at FROM categories WHERE parent_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*domain.Category{}
	for rows.Next() {
		var (
			id        string
			name      string
			parentID  *string
			createdAt time.Time
			deletedAt sql.NullTime
		)
		err := rows.Scan(&id, &name, &parentID, &createdAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		var deletedAtPtr *time.Time
		if deletedAt.Valid {
			deletedAtPtr = &deletedAt.Time
		}
		categories = append(categories, domain.NewCategoryFromRepository(id, name, parentID, createdAt, deletedAtPtr))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) FindRoots(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, deleted_at FROM categories WHERE parent_id IS NULL AND deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*domain.Category{}
	for rows.Next() {
		var (
			id        string
			name      string
			parentID  *string
			createdAt time.Time
			deletedAt sql.NullTime
		)
		err := rows.Scan(&id, &name, &parentID, &createdAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		var deletedAtPtr *time.Time
		if deletedAt.Valid {
			deletedAtPtr = &deletedAt.Time
		}
		categories = append(categories, domain.NewCategoryFromRepository(id, name, parentID, createdAt, deletedAtPtr))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `UPDATE categories SET name = $1, parent_id = $2 WHERE id = $3 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, category.Name(), category.ParentID(), category.ID())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

func (r *PostgresCategoryRepository) SoftDelete(ctx context.Context, categoryID string) error {
	query := `UPDATE categories SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), categoryID)
	return err
}
