package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresReviewRepository struct {
	db *sqlx.DB
}

func NewPostgresReviewRepository(db *sqlx.DB) *PostgresReviewRepository {
	return &PostgresReviewRepository{db: db}
}

const reviewColumns = `id, user_id, product_id, order_id, rating, comment, created_at`

type reviewRow struct {
	ID        uuid.UUID      `db:"id"`
	UserID    string         `db:"user_id"`
	ProductID string         `db:"product_id"`
	OrderID   string         `db:"order_id"`
	Rating    int            `db:"rating"`
	Comment   sql.NullString `db:"comment"`
	CreatedAt time.Time      `db:"created_at"`
}

// isUniqueViolation - reviewsning (user_id, product_id) unique cheklovi bo'yicha DB darajasida
// yakuniy himoya. SubmitReviewUseCase ExistsForUserAndProduct orqali oldindan tekshiradi,
// lekin bu constraint parallel so'rovlar (race condition) uchun so'nggi himoya sifatida qoladi.
func isUniqueViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" && pqErr.Constraint == constraint
	}
	return false
}

func (r *PostgresReviewRepository) Save(ctx context.Context, review *domain.Review) error {
	query := `INSERT INTO reviews (` + reviewColumns + `) VALUES (
		:id, :user_id, :product_id, :order_id, :rating, :comment, :created_at
	)`
	params := map[string]any{
		"id":         review.ID(),
		"user_id":    review.UserID(),
		"product_id": review.ProductID(),
		"order_id":   review.OrderID(),
		"rating":     review.Rating(),
		"comment":    nullIfEmpty(review.Comment()),
		"created_at": review.CreatedAt(),
	}

	if _, err := r.db.NamedExecContext(ctx, query, params); err != nil {
		if isUniqueViolation(err, "uq_reviews_user_product") {
			return domain.ErrReviewAlreadyExists
		}
		return err
	}
	return nil
}

func (r *PostgresReviewRepository) FindByProductID(ctx context.Context, productID string, page, pageSize int) ([]*domain.Review, int64, error) {
	countQuery := `SELECT COUNT(*) FROM reviews WHERE product_id = :product_id`
	countStmt, err := r.db.PrepareNamedContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countStmt.Close()

	var total int64
	if err := countStmt.GetContext(ctx, &total, map[string]any{"product_id": productID}); err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + reviewColumns + ` FROM reviews WHERE product_id = :product_id
		ORDER BY created_at DESC LIMIT :limit OFFSET :offset`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	params := map[string]any{
		"product_id": productID,
		"limit":      pageSize,
		"offset":     (page - 1) * pageSize,
	}

	var rows []reviewRow
	if err := stmt.SelectContext(ctx, &rows, params); err != nil {
		return nil, 0, err
	}

	return rowsToReviews(rows), total, nil
}

func (r *PostgresReviewRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Review, error) {
	query := `SELECT ` + reviewColumns + ` FROM reviews WHERE user_id = :user_id ORDER BY created_at DESC`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []reviewRow
	if err := stmt.SelectContext(ctx, &rows, map[string]any{"user_id": userID}); err != nil {
		return nil, err
	}

	return rowsToReviews(rows), nil
}

func (r *PostgresReviewRepository) ExistsForUserAndProduct(ctx context.Context, userID, productID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = :user_id AND product_id = :product_id)`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	params := map[string]any{"user_id": userID, "product_id": productID}

	var exists bool
	if err := stmt.GetContext(ctx, &exists, params); err != nil {
		return false, err
	}
	return exists, nil
}

func rowsToReviews(rows []reviewRow) []*domain.Review {
	reviews := make([]*domain.Review, 0, len(rows))
	for _, row := range rows {
		reviews = append(reviews, rowToReview(row))
	}
	return reviews
}

func rowToReview(row reviewRow) *domain.Review {
	comment := ""
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	return domain.NewReviewFromRepository(row.ID, row.UserID, row.ProductID, row.OrderID, row.Rating, comment, row.CreatedAt)
}

func nullIfEmpty(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
