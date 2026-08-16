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
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	ParentID  *string      `db:"parent_id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (r *PostgresCategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (id, name, parent_id, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, category.ID(), category.Name(), category.ParentID(), category.CreatedAt(), category.UpdatedAt(), category.DeletedAt())
	return err
}

func (r *PostgresCategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, updated_at, deleted_at FROM categories WHERE id = $1 AND deleted_at IS NULL`

	var (
		name      string
		parentID  *string
		createdAt time.Time
		updatedAt time.Time
		deletedAt *time.Time
	)

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&name, &parentID, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return domain.NewCategoryFromRepository(id, name, parentID, createdAt, updatedAt, deletedAt), nil
}

// FindAll returns all categories that are not deleted and have no parent.
func (r *PostgresCategoryRepository) FindAll(ctx context.Context, search string) ([]*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, updated_at, deleted_at
		          FROM categories
		          WHERE deleted_at IS NULL AND parent_id IS NULL AND name ILIKE :search
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
			row.ID, row.Name, row.ParentID, row.CreatedAt, row.UpdatedAt, deletedAtPtr,
		))
	}
	return categories, nil
}

func (r *PostgresCategoryRepository) FindByParentAndName(ctx context.Context, parentID *string, name string) (*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, updated_at, deleted_at
	          FROM categories
	          WHERE parent_id IS NOT DISTINCT FROM $1 AND name = $2 AND deleted_at IS NULL`

	var (
		id, catName string
		pID         *string
		createdAt   time.Time
		updatedAt   time.Time
		deletedAt   *time.Time
	)

	row := r.db.QueryRowContext(ctx, query, parentID, name)
	err := row.Scan(&id, &catName, &pID, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // topilmadi — bu xato emas, shunchaki mavjud emas
		}
		return nil, err
	}

	return domain.NewCategoryFromRepository(id, catName, pID, createdAt, updatedAt, deletedAt), nil
}

func (r *PostgresCategoryRepository) FindChildren(ctx context.Context, parentID string) ([]*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, updated_at, deleted_at FROM categories WHERE parent_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

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
			updatedAt time.Time
			deletedAt sql.NullTime
		)
		err := rows.Scan(&id, &name, &parentID, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		var deletedAtPtr *time.Time
		if deletedAt.Valid {
			deletedAtPtr = &deletedAt.Time
		}
		categories = append(categories, domain.NewCategoryFromRepository(id, name, parentID, createdAt, updatedAt, deletedAtPtr))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) FindRoots(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name, parent_id, created_at, updated_at, deleted_at FROM categories WHERE parent_id IS NULL AND deleted_at IS NULL ORDER BY created_at DESC`

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
			updatedAt time.Time
			deletedAt sql.NullTime
		)
		err := rows.Scan(&id, &name, &parentID, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		var deletedAtPtr *time.Time
		if deletedAt.Valid {
			deletedAtPtr = &deletedAt.Time
		}
		categories = append(categories, domain.NewCategoryFromRepository(id, name, parentID, createdAt, updatedAt, deletedAtPtr))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `UPDATE categories SET name = $1, parent_id = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, category.Name(), category.ParentID(), time.Now(), category.ID())

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
