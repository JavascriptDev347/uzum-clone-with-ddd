package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresGalleryPostRepository struct {
	db *sqlx.DB
}

func NewPostgresGalleryPostRepository(db *sqlx.DB) *PostgresGalleryPostRepository {
	return &PostgresGalleryPostRepository{db: db}
}

const galleryPostColumns = `id, images, description, created_at`

type galleryPostRow struct {
	ID          uuid.UUID `db:"id"`
	Images      []byte    `db:"images"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

// Save - images'ni catalog'ning Product.Images'i bilan bir xil naqshda (JSONB ustun)
// saqlaydi, izchillik uchun.
func (r *PostgresGalleryPostRepository) Save(ctx context.Context, post *domain.GalleryPost) error {
	images := post.Images()
	if images == nil {
		images = []domain.GalleryImage{}
	}
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return err
	}

	query := `INSERT INTO gallery_posts (` + galleryPostColumns + `) VALUES (
		:id, :images, :description, :created_at
	)`
	params := map[string]any{
		"id":          post.ID(),
		"images":      string(imagesJSON),
		"description": post.Description(),
		"created_at":  post.CreatedAt(),
	}

	_, err = r.db.NamedExecContext(ctx, query, params)
	return err
}

func (r *PostgresGalleryPostRepository) FindAll(ctx context.Context, page, pageSize int) ([]*domain.GalleryPost, int64, error) {
	var total int64
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM gallery_posts`); err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + galleryPostColumns + ` FROM gallery_posts
		ORDER BY created_at DESC LIMIT :limit OFFSET :offset`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	params := map[string]any{
		"limit":  pageSize,
		"offset": (page - 1) * pageSize,
	}

	var rows []galleryPostRow
	if err := stmt.SelectContext(ctx, &rows, params); err != nil {
		return nil, 0, err
	}

	posts := make([]*domain.GalleryPost, 0, len(rows))
	for _, row := range rows {
		post, err := rowToGalleryPost(row)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}
	return posts, total, nil
}

// Delete - postni o'chiradi va DELETE ... RETURNING orqali o'chirilgan qatorni qaytaradi -
// chaqiruvchi (DeleteGalleryPostUseCase) shu orqali alohida FindByID so'rovisiz rasmlarning
// PublicID'larini olib, ularni S3'dan tozalay oladi.
func (r *PostgresGalleryPostRepository) Delete(ctx context.Context, id string) (*domain.GalleryPost, error) {
	query := `DELETE FROM gallery_posts WHERE id = :id RETURNING ` + galleryPostColumns
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var row galleryPostRow
	if err := stmt.GetContext(ctx, &row, map[string]any{"id": id}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrGalleryPostNotFound
		}
		return nil, err
	}

	return rowToGalleryPost(row)
}

func rowToGalleryPost(row galleryPostRow) (*domain.GalleryPost, error) {
	images := []domain.GalleryImage{}
	if len(row.Images) > 0 {
		if err := json.Unmarshal(row.Images, &images); err != nil {
			return nil, err
		}
	}
	return domain.NewGalleryPostFromRepository(row.ID, images, row.Description, row.CreatedAt), nil
}
