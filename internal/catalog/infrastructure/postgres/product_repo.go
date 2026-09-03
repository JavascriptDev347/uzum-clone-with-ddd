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

const productColumns = `id, name_uz, name_eng, name_ru, description_uz, description_eng, description_ru, images,
	category_id, price_amount, price_currency, discount_price_amount, slug, is_available,
	rating, stock, sold_count, tag_uz, tag_eng, tag_ru,
	created_at, updated_at, deleted_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProduct(s rowScanner) (*domain.Product, error) {
	var (
		id             string
		nameUz         string
		nameEng        string
		nameRu         string
		descriptionUz  string
		descriptionEng string
		descriptionRu  string
		imagesRaw      []byte
		categoryID     string
		priceAmount    int64
		priceCurrency  string
		discountAmount sql.NullInt64
		slug           string
		isAvailable    bool
		rating         float64
		stock          int
		soldCount      int
		tagUz          sql.NullString
		tagEng         sql.NullString
		tagRu          sql.NullString
		createdAt      time.Time
		updatedAt      time.Time
		deletedAt      sql.NullTime
	)

	err := s.Scan(
		&id, &nameUz, &nameEng, &nameRu, &descriptionUz, &descriptionEng, &descriptionRu, &imagesRaw,
		&categoryID, &priceAmount, &priceCurrency, &discountAmount, &slug, &isAvailable,
		&rating, &stock, &soldCount, &tagUz, &tagEng, &tagRu,
		&createdAt, &updatedAt, &deletedAt,
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

	var tagUzPtr, tagEngPtr, tagRuPtr *string
	if tagUz.Valid {
		v := tagUz.String
		tagUzPtr = &v
	}
	if tagEng.Valid {
		v := tagEng.String
		tagEngPtr = &v
	}
	if tagRu.Valid {
		v := tagRu.String
		tagRuPtr = &v
	}

	var deletedAtPtr *time.Time
	if deletedAt.Valid {
		deletedAtPtr = &deletedAt.Time
	}

	return domain.NewProductFromRepository(domain.ProductFromRepositoryParams{
		ID:             id,
		NameUz:         nameUz,
		NameEng:        nameEng,
		NameRu:         nameRu,
		DescriptionUz:  descriptionUz,
		DescriptionEng: descriptionEng,
		DescriptionRu:  descriptionRu,
		Images:         images,
		CategoryID:     categoryID,
		Price:          price,
		DiscountPrice:  discountPrice,
		Slug:           slug,
		IsAvailable:    isAvailable,
		Rating:         rating,
		Stock:          stock,
		SoldCount:      soldCount,
		TagUz:          tagUzPtr,
		TagEng:         tagEngPtr,
		TagRu:          tagRuPtr,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DeletedAt:      deletedAtPtr,
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
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23
	)`

	_, err = r.db.ExecContext(ctx, query,
		p.ID(), p.NameUz(), p.NameEng(), p.NameRu(), p.DescriptionUz(), p.DescriptionEng(), p.DescriptionRu(), string(imagesJSON),
		p.CategoryID(), p.Price().Amount(), p.Price().Currency(), discountAmount, p.Slug(), p.IsAvailable(),
		p.Rating(), p.Stock(), p.SoldCount(), p.TagUz(), p.TagEng(), p.TagRu(),
		p.CreatedAt(), p.UpdatedAt(), p.DeletedAt(),
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

func (r *PostgresProductRepository) findAll(ctx context.Context, search, categoryID string, includeDeleted bool, page, pageSize int) ([]*domain.Product, int64, error) {
	where := ` WHERE (name_uz ILIKE $1 OR name_eng ILIKE $1 OR name_ru ILIKE $1)`
	args := []any{"%" + search + "%"}

	if categoryID != "" {
		args = append(args, categoryID)
		where += ` AND category_id = $` + strconv.Itoa(len(args))
	}
	if !includeDeleted {
		where += ` AND deleted_at IS NULL`
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM products` + where
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	query := `SELECT ` + productColumns + ` FROM products` + where +
		` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)-1) + ` OFFSET $` + strconv.Itoa(len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []*domain.Product{}
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *PostgresProductRepository) FindAll(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*domain.Product, int64, error) {
	return r.findAll(ctx, search, categoryID, false, page, pageSize)
}

func (r *PostgresProductRepository) FindAllIncludingDeleted(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*domain.Product, int64, error) {
	return r.findAll(ctx, search, categoryID, true, page, pageSize)
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
		name_uz=$1, name_eng=$2, name_ru=$3, description_uz=$4, description_eng=$5, description_ru=$6, images=$7,
		category_id=$8, price_amount=$9, price_currency=$10, discount_price_amount=$11, slug=$12, is_available=$13,
		rating=$14, stock=$15, sold_count=$16, tag_uz=$17, tag_eng=$18, tag_ru=$19, updated_at=$20
		WHERE id=$21 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		p.NameUz(), p.NameEng(), p.NameRu(), p.DescriptionUz(), p.DescriptionEng(), p.DescriptionRu(), string(imagesJSON),
		p.CategoryID(), p.Price().Amount(), p.Price().Currency(), discountAmount, p.Slug(), p.IsAvailable(),
		p.Rating(), p.Stock(), p.SoldCount(), p.TagUz(), p.TagEng(), p.TagRu(), p.UpdatedAt(),
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
