package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresProductRepository struct {
	db *sqlx.DB
}

func NewPostgresProductRepository(db *sqlx.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Save(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (id, name, category_id, price_amount, price_currency, image_url, image_public_id, created_at, deleted_at)
	VALUES (:id, :name, :category_id, :price_amount, :price_currency, :image_url, :image_public_id, :created_at, :deleted_at)`

	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":              product.ID(),
		"name":            product.Name(),
		"category_id":     product.CategoryID(),
		"price_amount":    product.Price().Amount(),
		"price_currency":  product.Price().Currency(),
		"image_url":       product.ImageURL(),
		"image_public_id": product.ImagePublicID(),
		"created_at":      product.CreatedAt(),
		"deleted_at":      product.DeletedAt(),
	})
	return err
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `SELECT id, name, category_id, price_amount, price_currency, image_url, image_public_id, created_at, deleted_at
			FROM products
			WHERE id = $1 AND deleted_at IS NULL`

	var (
		productID     string
		name          string
		categoryID    string
		priceAmount   int64
		priceCurrency string
		imageURL      sql.NullString
		imagePublicID sql.NullString
		createdAt     time.Time
		deletedAt     sql.NullTime
	)

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&productID, &name, &categoryID, &priceAmount, &priceCurrency, &imageURL, &imagePublicID, &createdAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	price, err := domain.NewMoney(priceAmount, priceCurrency)
	if err != nil {
		return nil, err
	}

	var deletedAtPtr *time.Time
	if deletedAt.Valid {
		deletedAtPtr = &deletedAt.Time
	}

	return domain.NewProductFromRepository(
		productID, name, categoryID, price,
		imageURL.String, imagePublicID.String,
		createdAt, deletedAtPtr,
	), nil
}

func (r *PostgresProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	query := `SELECT id, name, category_id, price_amount, price_currency, image_url, image_public_id, created_at, deleted_at
			FROM products
			WHERE deleted_at IS NULL
			ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*domain.Product{}
	for rows.Next() {
		var (
			productID     string
			name          string
			categoryID    string
			priceAmount   int64
			priceCurrency string
			imageURL      sql.NullString
			imagePublicID sql.NullString
			createdAt     time.Time
			deletedAt     sql.NullTime
		)

		err := rows.Scan(&productID, &name, &categoryID, &priceAmount, &priceCurrency, &imageURL, &imagePublicID, &createdAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		price, err := domain.NewMoney(priceAmount, priceCurrency)
		if err != nil {
			return nil, err
		}

		var deletedAtPtr *time.Time
		if deletedAt.Valid {
			deletedAtPtr = &deletedAt.Time
		}

		products = append(products, domain.NewProductFromRepository(
			productID, name, categoryID, price,
			imageURL.String, imagePublicID.String,
			createdAt, deletedAtPtr,
		))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products
			SET name=:name, category_id=:category_id, price_amount=:price_amount, price_currency=:price_currency,
				image_url=:image_url, image_public_id=:image_public_id
			WHERE id=:id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":              product.ID(),
		"name":            product.Name(),
		"category_id":     product.CategoryID(),
		"price_amount":    product.Price().Amount(),
		"price_currency":  product.Price().Currency(),
		"image_url":       product.ImageURL(),
		"image_public_id": product.ImagePublicID(),
	})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *PostgresProductRepository) SoftDelete(ctx context.Context, productID string) error {
	query := `UPDATE products SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), productID)
	return err
}
