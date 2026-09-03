package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresEventRepository struct {
	db *sqlx.DB
}

func NewPostgresEventRepository(db *sqlx.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

type eventRow struct {
	ID string `db:"id"`

	EyebrowUz  string `db:"eyebrow_uz"`
	EyebrowEng string `db:"eyebrow_eng"`
	EyebrowRu  string `db:"eyebrow_ru"`

	TitleUz  string `db:"title_uz"`
	TitleEng string `db:"title_eng"`
	TitleRu  string `db:"title_ru"`

	SubtitleUz  string `db:"subtitle_uz"`
	SubtitleEng string `db:"subtitle_eng"`
	SubtitleRu  string `db:"subtitle_ru"`

	CTAUz  string `db:"cta_uz"`
	CTAEng string `db:"cta_eng"`
	CTARu  string `db:"cta_ru"`

	ImageURL      string       `db:"image_url"`
	ImagePublicID string       `db:"image_public_id"`
	CategoryID    string       `db:"category_id"`
	IsRoot        bool         `db:"is_root"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
}

func (row eventRow) toDomain() *domain.Event {
	var deletedAtPtr *time.Time
	if row.DeletedAt.Valid {
		deletedAtPtr = &row.DeletedAt.Time
	}
	return domain.NewEventFromRepository(domain.EventFromRepositoryParams{
		ID:            row.ID,
		EyebrowUz:     row.EyebrowUz,
		EyebrowEng:    row.EyebrowEng,
		EyebrowRu:     row.EyebrowRu,
		TitleUz:       row.TitleUz,
		TitleEng:      row.TitleEng,
		TitleRu:       row.TitleRu,
		SubtitleUz:    row.SubtitleUz,
		SubtitleEng:   row.SubtitleEng,
		SubtitleRu:    row.SubtitleRu,
		CTAUz:         row.CTAUz,
		CTAEng:        row.CTAEng,
		CTARu:         row.CTARu,
		ImageURL:      row.ImageURL,
		ImagePublicID: row.ImagePublicID,
		CategoryID:    row.CategoryID,
		IsRoot:        row.IsRoot,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		DeletedAt:     deletedAtPtr,
	})
}

const eventColumns = `id, eyebrow_uz, eyebrow_eng, eyebrow_ru, title_uz, title_eng, title_ru,
	subtitle_uz, subtitle_eng, subtitle_ru, cta_uz, cta_eng, cta_ru,
	image_url, image_public_id, category_id, is_root, created_at, updated_at, deleted_at`

func (r *PostgresEventRepository) Save(ctx context.Context, event *domain.Event) error {
	query := `INSERT INTO events (` + eventColumns + `)
		VALUES (:id, :eyebrow_uz, :eyebrow_eng, :eyebrow_ru, :title_uz, :title_eng, :title_ru,
			:subtitle_uz, :subtitle_eng, :subtitle_ru, :cta_uz, :cta_eng, :cta_ru,
			:image_url, :image_public_id, :category_id, :is_root, :created_at, :updated_at, :deleted_at)`

	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":              event.ID(),
		"eyebrow_uz":      event.EyebrowUz(),
		"eyebrow_eng":     event.EyebrowEng(),
		"eyebrow_ru":      event.EyebrowRu(),
		"title_uz":        event.TitleUz(),
		"title_eng":       event.TitleEng(),
		"title_ru":        event.TitleRu(),
		"subtitle_uz":     event.SubtitleUz(),
		"subtitle_eng":    event.SubtitleEng(),
		"subtitle_ru":     event.SubtitleRu(),
		"cta_uz":          event.CTAUz(),
		"cta_eng":         event.CTAEng(),
		"cta_ru":          event.CTARu(),
		"image_url":       event.ImageURL(),
		"image_public_id": event.ImagePublicID(),
		"category_id":     event.CategoryID(),
		"is_root":         event.IsRoot(),
		"created_at":      event.CreatedAt(),
		"updated_at":      event.UpdatedAt(),
		"deleted_at":      event.DeletedAt(),
	})
	return err
}

func (r *PostgresEventRepository) FindByID(ctx context.Context, id string) (*domain.Event, error) {
	query := `SELECT ` + eventColumns + ` FROM events WHERE id = $1 AND deleted_at IS NULL`

	var row eventRow
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		return nil, err
	}
	return row.toDomain(), nil
}

func (r *PostgresEventRepository) FindAll(ctx context.Context) ([]*domain.Event, error) {
	query := `SELECT ` + eventColumns + ` FROM events WHERE deleted_at IS NULL ORDER BY is_root DESC, created_at DESC`

	var rows []eventRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	events := make([]*domain.Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, row.toDomain())
	}
	return events, nil
}

func (r *PostgresEventRepository) FindAllIncludingDeleted(ctx context.Context) ([]*domain.Event, error) {
	query := `SELECT ` + eventColumns + ` FROM events ORDER BY is_root DESC, created_at DESC`

	var rows []eventRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	events := make([]*domain.Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, row.toDomain())
	}
	return events, nil
}

func (r *PostgresEventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `UPDATE events
		SET eyebrow_uz=:eyebrow_uz, eyebrow_eng=:eyebrow_eng, eyebrow_ru=:eyebrow_ru,
			title_uz=:title_uz, title_eng=:title_eng, title_ru=:title_ru,
			subtitle_uz=:subtitle_uz, subtitle_eng=:subtitle_eng, subtitle_ru=:subtitle_ru,
			cta_uz=:cta_uz, cta_eng=:cta_eng, cta_ru=:cta_ru,
			category_id=:category_id, is_root=:is_root, updated_at=:updated_at
		WHERE id=:id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"eyebrow_uz":   event.EyebrowUz(),
		"eyebrow_eng":  event.EyebrowEng(),
		"eyebrow_ru":   event.EyebrowRu(),
		"title_uz":     event.TitleUz(),
		"title_eng":    event.TitleEng(),
		"title_ru":     event.TitleRu(),
		"subtitle_uz":  event.SubtitleUz(),
		"subtitle_eng": event.SubtitleEng(),
		"subtitle_ru":  event.SubtitleRu(),
		"cta_uz":       event.CTAUz(),
		"cta_eng":      event.CTAEng(),
		"cta_ru":       event.CTARu(),
		"category_id":  event.CategoryID(),
		"is_root":      event.IsRoot(),
		"updated_at":   time.Now(),
		"id":           event.ID(),
	})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}
	return nil
}

func (r *PostgresEventRepository) SoftDelete(ctx context.Context, id string) error {
	query := `UPDATE events SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
