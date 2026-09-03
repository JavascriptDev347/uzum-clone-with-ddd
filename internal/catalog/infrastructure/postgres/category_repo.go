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
	NameUz        string       `db:"name_uz"`
	NameEng       string       `db:"name_eng"`
	NameRu        string       `db:"name_ru"`
	ImageURL      string       `db:"image_url"`
	ImagePublicID string       `db:"image_public_id"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
}

func (r *PostgresCategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (id, name_uz, name_eng, name_ru, image_url, image_public_id, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query, category.ID(), category.NameUz(), category.NameEng(), category.NameRu(), category.ImageURL(), category.ImagePublicID(), category.CreatedAt(), category.UpdatedAt(), category.DeletedAt())
	return err
}

func (r *PostgresCategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `SELECT id, name_uz, name_eng, name_ru, image_url, image_public_id, created_at, updated_at, deleted_at FROM categories WHERE id = $1 AND deleted_at IS NULL`

	var (
		categoryID    string
		nameUz        string
		nameEng       string
		nameRu        string
		imageURL      string
		imagePublicID string
		createdAt     time.Time
		updatedAt     time.Time
		deletedAt     *time.Time
	)

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&categoryID, &nameUz, &nameEng, &nameRu, &imageURL, &imagePublicID, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return domain.NewCategoryFromRepository(id, nameUz, nameEng, nameRu, imageURL, imagePublicID, createdAt, updatedAt, deletedAt), nil
}

func (r *PostgresCategoryRepository) FindAll(ctx context.Context, search string) ([]*domain.Category, error) {
	query := `SELECT id, name_uz, name_eng, name_ru, image_url, image_public_id, created_at, updated_at, deleted_at
		          FROM categories
		          WHERE deleted_at IS NULL AND (name_uz ILIKE :search OR name_eng ILIKE :search OR name_ru ILIKE :search)
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
			row.ID, row.NameUz, row.NameEng, row.NameRu, row.ImageURL, row.ImagePublicID, row.CreatedAt, row.UpdatedAt, deletedAtPtr,
		))
	}
	return categories, nil
}

func (r *PostgresCategoryRepository) FindAllIncludingDeleted(ctx context.Context, search string) ([]*domain.Category, error) {
	query := `SELECT id, name_uz, name_eng, name_ru, image_url, image_public_id, created_at, updated_at, deleted_at
	          FROM categories
	          WHERE (name_uz ILIKE :search OR name_eng ILIKE :search OR name_ru ILIKE :search)
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
			row.ID, row.NameUz, row.NameEng, row.NameRu, row.ImageURL, row.ImagePublicID, row.CreatedAt, row.UpdatedAt, deletedAtPtr,
		))
	}
	return categories, nil
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
			UPDATE categories
			SET name_uz = :name_uz, name_eng = :name_eng, name_ru = :name_ru, updated_at = :updated_at
			WHERE id = :id AND deleted_at IS NULL
		`
	params := map[string]any{
		"name_uz":    category.NameUz(),
		"name_eng":   category.NameEng(),
		"name_ru":    category.NameRu(),
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
