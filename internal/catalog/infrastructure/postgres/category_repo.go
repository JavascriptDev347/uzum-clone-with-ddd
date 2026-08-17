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

type categoryRow struct {
	ID            string       `db:"id"`
	Name          string       `db:"name"`
	ImageURL      string       `db:"image_url"`
	ImagePublicID string       `db:"image_public_id"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
}

func (r *PostgresCategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (id, name, image_url, image_public_id, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query, category.ID(), category.Name(), category.ImageURL(), category.ImagePublicID(), category.CreatedAt(), category.UpdatedAt(), category.DeletedAt())
	return err
}

func (r *PostgresCategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `SELECT id, name, image_url, image_public_id, created_at, updated_at, deleted_at FROM categories WHERE id = $1 AND deleted_at IS NULL`

	var (
		categoryID    string
		name          string
		imageURL      string
		imagePublicID string
		createdAt     time.Time
		updatedAt     time.Time
		deletedAt     *time.Time
	)

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&categoryID, &name, &imageURL, &imagePublicID, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return domain.NewCategoryFromRepository(id, name, imageURL, imagePublicID, createdAt, updatedAt, deletedAt), nil
}

func (r *PostgresCategoryRepository) FindAll(ctx context.Context, search string) ([]*domain.Category, error) {
	query := `SELECT id, name, image_url, image_public_id, created_at, updated_at, deleted_at
		          FROM categories
		          WHERE deleted_at IS NULL AND name ILIKE :search
		          ORDER BY created_at DESC`

	params := map[string]any{
		"search": "%" + search + "%",
	}

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []categoryRow
	if err := stmt.SelectContext(ctx, &rows, params); err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, 0, len(rows))
	for _, row := range rows {
		var deletedAtPtr *time.Time
		if row.DeletedAt.Valid {
			deletedAtPtr = &row.DeletedAt.Time
		}
		categories = append(categories, domain.NewCategoryFromRepository(
			row.ID, row.Name, row.ImageURL, row.ImagePublicID, row.CreatedAt, row.UpdatedAt, deletedAtPtr,
		))
	}
	return categories, nil
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
			UPDATE categories
			SET name = :name, updated_at = :updated_at
			WHERE id = :id AND deleted_at IS NULL
		`
	params := map[string]interface{}{
		"name":       category.Name(),
		"updated_at": time.Now(),
		"id":         category.ID(),
	}
	result, err := r.db.NamedExecContext(ctx, query, params)
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
	query := `UPDATE categories SET deleted_at=$1  WHERE id=$2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), categoryID)
	return err
}
