package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresWishlistRepository struct {
	db *sql.DB
}

func NewPostgresWishlistRepository(db *sql.DB) *PostgresWishlistRepository {
	return &PostgresWishlistRepository{db: db}
}

type wishlistRow struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	ProductID string       `db:"product_id"`
	CreatedAt sql.NullTime `db:"created_at"`
}

func (r *PostgresWishlistRepository) Save(ctx context.Context, wishlist *domain.Wishlist) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // agar Commit muvaffaqiyatli bo'lsa, bu no-op bo'ladi

	upsertQuery := `
		INSERT INTO wishlists (id, user_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING`
	if _, err := tx.ExecContext(ctx, upsertQuery, wishlist.ID(), wishlist.UserID(), wishlist.CreatedAt()); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM wishlist_items WHERE wishlist_id = $1`, wishlist.ID()); err != nil {
		return err
	}

	insertItemQuery := `INSERT INTO wishlist_items (id, wishlist_id, product_id, added_at) VALUES ($1, $2, $3, $4)`
	for _, item := range wishlist.Items() {
		if _, err := tx.ExecContext(ctx, insertItemQuery, uuid.New().String(), wishlist.ID(), item.ProductID(), item.AddedAt()); err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return domain.ErrProductAlreadyInWishlist
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresWishlistRepository) FindByUserID(ctx context.Context, userID string) (*domain.Wishlist, error) {
	query := `
			SELECT w.id, w.user_id, w.created_at, wi.product_id, wi.added_at
			FROM wishlists w
			LEFT JOIN wishlist_items wi ON wi.wishlist_id = w.id
			WHERE w.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		wishlistID string
		createdAt  time.Time
		items      []domain.WishlistItem
		found      bool
	)

	for rows.Next() {
		found = true
		var productID sql.NullString
		var addedAt sql.NullTime

		if err := rows.Scan(&wishlistID, &userID, &createdAt, &productID, &addedAt); err != nil {
			return nil, err
		}

		// LEFT JOIN natijasida bo'sh wishlist uchun productID NULL bo'lishi mumkin
		if productID.Valid {
			items = append(items, domain.NewWishlistItemFromRepository(productID.String, addedAt.Time))

		}

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !found {
		return nil, domain.ErrWishlistNotFound
	}

	return domain.NewWishlistFromRepository(wishlistID, userID, items, createdAt), nil

}
