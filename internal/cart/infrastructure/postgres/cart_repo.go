package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
	"github.com/google/uuid"
)

type PostgresCartRepository struct {
	db *sql.DB
}

func NewPostgresCartRepository(db *sql.DB) *PostgresCartRepository {
	return &PostgresCartRepository{db: db}
}

type cartRow struct {
	ID        uuid.UUID `db:"id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type cartItemRow struct {
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
}

func (r *PostgresCartRepository) Save(ctx context.Context, c *domain.Cart) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			INSERT INTO carts (id, user_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET updated_at = EXCLUDED.updated_at
		`, c.ID(), c.UserID(), c.CreatedAt(), c.UpdatedAt())
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, c.ID()); err != nil {
		return err
	}

	for _, item := range c.Items() {
		_, err = tx.ExecContext(ctx, `
				INSERT INTO cart_items (cart_id, product_id, quantity)
				VALUES ($1, $2, $3)
			`, c.ID(), item.ProductID(), item.Quantity())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresCartRepository) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {

	var (
		id        uuid.UUID
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, `
			SELECT id, created_at, updated_at FROM carts WHERE user_id = $1
		`, userID).Scan(&id, &createdAt, &updatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // topilmadi — bu xato emas, use case lazy-create qiladi
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
			SELECT product_id, quantity FROM cart_items WHERE cart_id = $1
		`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.CartItem, 0)
	for rows.Next() {
		var (
			productID string
			quantity  int
		)
		if err := rows.Scan(&productID, &quantity); err != nil {
			return nil, err
		}
		item, err := domain.NewCartItem(productID, quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return domain.NewCartFromRepository(id, userID, items, createdAt, updatedAt), nil
}
