package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresProductRepository struct {
	db *sqlx.DB
}

func NewPostgresProductRepository(db *sqlx.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

const productColumns = `id, name, description, images, video_url_youtube, video_url_instagram, category_id,
	price_amount, price_currency, discount_price_amount, slug, is_available,
	rating, stock, sold_count, flower_types, color, stem_count,
	packaging_type, freshness_lifespan, care_instructions, occasions,
	allow_custom_card, compatible_addons, created_at, updated_at, deleted_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProduct(s rowScanner) (*domain.Product, error) {
	var (
		id                string
		name              string
		description       string
		imagesRaw         []byte
		videoURLYoutube   string
		videoURLInstagram string
		categoryID        string
		priceAmount       int64
		priceCurrency     string
		discountAmount    sql.NullInt64
		slug              string
		isAvailable       bool
		rating            float64
		stock             int
		soldCount         int
		flowerTypes       []string
		color             string
		stemCount         int
		packagingType     string
		freshnessLifespan int
		careInstructions  sql.NullString
		occasions         []string
		allowCustomCard   sql.NullBool
		compatibleAddons  []string
		createdAt         time.Time
		updatedAt         time.Time
		deletedAt         sql.NullTime
	)

	err := s.Scan(
		&id, &name, &description, &imagesRaw, &videoURLYoutube, &videoURLInstagram, &categoryID,
		&priceAmount, &priceCurrency, &discountAmount, &slug, &isAvailable,
		&rating, &stock, &soldCount, pq.Array(&flowerTypes), &color, &stemCount,
		&packagingType, &freshnessLifespan, &careInstructions, pq.Array(&occasions),
		&allowCustomCard, pq.Array(&compatibleAddons), &createdAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	price, err := domain.NewMoney(priceAmount, priceCurrency)
	if err != nil {
		return nil, err
	}

	var discountPrice *domain.Money
	if discountAmount.Valid {
		dp, err := domain.NewMoney(discountAmount.Int64, priceCurrency)
		if err != nil {
			return nil, err
		}
		discountPrice = &dp
	}

	images := []domain.ProductImage{}
	if len(imagesRaw) > 0 {
		if err := json.Unmarshal(imagesRaw, &images); err != nil {
			return nil, err
		}
	}

	var careInstructionsPtr *string
	if careInstructions.Valid {
		v := careInstructions.String
		careInstructionsPtr = &v
	}

	var allowCustomCardPtr *bool
	if allowCustomCard.Valid {
		v := allowCustomCard.Bool
		allowCustomCardPtr = &v
	}

	var deletedAtPtr *time.Time
	if deletedAt.Valid {
		deletedAtPtr = &deletedAt.Time
	}

	return domain.NewProductFromRepository(domain.ProductFromRepositoryParams{
		ID:                id,
		Name:              name,
		Description:       description,
		Images:            images,
		VideoURLYoutube:   videoURLYoutube,
		VideoURLInstagram: videoURLInstagram,
		CategoryID:        categoryID,
		Price:             price,
		DiscountPrice:     discountPrice,
		Slug:              slug,
		IsAvailable:       isAvailable,
		Rating:            rating,
		Stock:             stock,
		SoldCount:         soldCount,
		FlowerTypes:       flowerTypes,
		Color:             color,
		StemCount:         stemCount,
		PackagingType:     domain.PackagingType(packagingType),
		FreshnessLifespan: domain.FreshnessLifespan(freshnessLifespan),
		CareInstructions:  careInstructionsPtr,
		Occasions:         occasions,
		AllowCustomCard:   allowCustomCardPtr,
		CompatibleAddons:  compatibleAddons,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		DeletedAt:         deletedAtPtr,
	}), nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" && pqErr.Constraint == constraint
	}
	return false
}

func (r *PostgresProductRepository) Save(ctx context.Context, p *domain.Product) error {
	images := p.Images()
	if images == nil {
		images = []domain.ProductImage{}
	}
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return err
	}

	var discountAmount *int64
	if dp := p.DiscountPrice(); dp != nil {
		v := dp.Amount()
		discountAmount = &v
	}

	query := `INSERT INTO products (` + productColumns + `) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27
	)`

	_, err = r.db.ExecContext(ctx, query,
		p.ID(), p.Name(), p.Description(), string(imagesJSON), p.VideoURLYoutube(), p.VideoURLInstagram(), p.CategoryID(),
		p.Price().Amount(), p.Price().Currency(), discountAmount, p.Slug(), p.IsAvailable(),
		p.Rating(), p.Stock(), p.SoldCount(), pq.Array(p.FlowerTypes()), p.Color(), p.StemCount(),
		string(p.PackagingType()), int(p.FreshnessLifespan()), p.CareInstructions(), pq.Array(p.Occasions()),
		p.AllowCustomCard(), pq.Array(p.CompatibleAddons()), p.CreatedAt(), p.UpdatedAt(), p.DeletedAt(),
	)
	if isUniqueViolation(err, "products_slug_key") {
		return domain.ErrProductSlugTaken
	}
	return err
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE id = $1 AND deleted_at IS NULL`

	product, err := scanProduct(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *PostgresProductRepository) FindBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE slug = $1 AND deleted_at IS NULL`

	product, err := scanProduct(r.db.QueryRowContext(ctx, query, slug))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *PostgresProductRepository) findAll(ctx context.Context, search, categoryID string, includeDeleted bool) ([]*domain.Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE name ILIKE $1`
	args := []any{"%" + search + "%"}

	if categoryID != "" {
		args = append(args, categoryID)
		query += ` AND category_id = $` + strconv.Itoa(len(args))
	}
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*domain.Product{}
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) FindAll(ctx context.Context, search string, categoryID string) ([]*domain.Product, error) {
	return r.findAll(ctx, search, categoryID, false)
}

func (r *PostgresProductRepository) FindAllIncludingDeleted(ctx context.Context, search string, categoryID string) ([]*domain.Product, error) {
	return r.findAll(ctx, search, categoryID, true)
}

func (r *PostgresProductRepository) Update(ctx context.Context, p *domain.Product) error {
	images := p.Images()
	if images == nil {
		images = []domain.ProductImage{}
	}
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return err
	}

	var discountAmount *int64
	if dp := p.DiscountPrice(); dp != nil {
		v := dp.Amount()
		discountAmount = &v
	}

	query := `UPDATE products SET
		name=$1, description=$2, images=$3, video_url_youtube=$4, video_url_instagram=$5, category_id=$6,
		price_amount=$7, price_currency=$8, discount_price_amount=$9, slug=$10, is_available=$11,
		rating=$12, stock=$13, sold_count=$14, flower_types=$15, color=$16, stem_count=$17,
		packaging_type=$18, freshness_lifespan=$19, care_instructions=$20, occasions=$21,
		allow_custom_card=$22, compatible_addons=$23, updated_at=$24
		WHERE id=$25 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		p.Name(), p.Description(), string(imagesJSON), p.VideoURLYoutube(), p.VideoURLInstagram(), p.CategoryID(),
		p.Price().Amount(), p.Price().Currency(), discountAmount, p.Slug(), p.IsAvailable(),
		p.Rating(), p.Stock(), p.SoldCount(), pq.Array(p.FlowerTypes()), p.Color(), p.StemCount(),
		string(p.PackagingType()), int(p.FreshnessLifespan()), p.CareInstructions(), pq.Array(p.Occasions()),
		p.AllowCustomCard(), pq.Array(p.CompatibleAddons()), p.UpdatedAt(),
		p.ID(),
	)
	if isUniqueViolation(err, "products_slug_key") {
		return domain.ErrProductSlugTaken
	}
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
