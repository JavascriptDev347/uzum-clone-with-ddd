package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyProductID       = errors.New("catalog: product ID bo'sh bo'lishi mumkin emas")
	ErrEmptyProductName     = errors.New("catalog: product nomi (uz, eng, ru) bo'sh bo'lishi mumkin emas")
	ErrEmptyProductSlug     = errors.New("catalog: product slug bo'sh bo'lishi mumkin emas")
	ErrEmptyCategoryID      = errors.New("catalog: category ID bo'sh bo'lishi mumkin emas")
	ErrTooManyProductImages = errors.New("catalog: mahsulot uchun eng ko'pi bilan 5 ta rasm yuklash mumkin")
	ErrInvalidRating        = errors.New("catalog: reyting 1 dan 5 gacha bo'lishi kerak")
	ErrNegativeStock        = errors.New("catalog: stock manfiy bo'lishi mumkin emas")
	ErrNegativeSoldCount    = errors.New("catalog: sotilgan mahsulotlar soni manfiy bo'lishi mumkin emas")
	ErrDiscountTooHigh      = errors.New("catalog: chegirma narxi asosiy narxdan kichik va bir xil valyutada bo'lishi kerak")
	ErrProductSlugTaken     = errors.New("catalog: bu slug allaqachon band")
)

const (
	// MaxProductImages - mahsulot uchun yuklash mumkin bo'lgan eng ko'p rasmlar soni.
	MaxProductImages = 5
	// DefaultRating - mahsulot yaratilganda reyting ko'rsatilmasa qo'llaniladigan qiymat.
	DefaultRating = 1
)

// ProductImage - mahsulotning bitta rasmi (Cloudinary'dagi manzili va public id'si).
type ProductImage struct {
	URL      string
	PublicID string
}

type Product struct {
	id string

	nameUz  string
	nameEng string
	nameRu  string

	descriptionUz  string
	descriptionEng string
	descriptionRu  string

	images     []ProductImage
	categoryID string

	price         Money
	discountPrice *Money // ixtiyoriy, bo'lsa umumiy summa shu bo'yicha hisoblanadi

	slug        string
	isAvailable bool
	rating      float64 // 1-5, default 1

	stock     int
	soldCount int

	// tagUz/tagEng/tagRu - ixtiyoriy belgi (badge), masalan "bestseller"
	tagUz  *string
	tagEng *string
	tagRu  *string

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

// NewProductParams - yangi mahsulot yaratish uchun kerakli ma'lumotlar.
type NewProductParams struct {
	ID             string
	NameUz         string
	NameEng        string
	NameRu         string
	DescriptionUz  string
	DescriptionEng string
	DescriptionRu  string
	Images         []ProductImage
	CategoryID     string
	Price          Money
	DiscountPrice  *Money
	Slug           string
	IsAvailable    bool
	Rating         float64
	Stock          int
	TagUz          *string
	TagEng         *string
	TagRu          *string
}

func NewProduct(p NewProductParams) (*Product, error) {
	if p.ID == "" {
		return nil, ErrEmptyProductID
	}
	if p.NameUz == "" || p.NameEng == "" || p.NameRu == "" {
		return nil, ErrEmptyProductName
	}
	if p.CategoryID == "" {
		return nil, ErrEmptyCategoryID
	}
	if p.Slug == "" {
		return nil, ErrEmptyProductSlug
	}
	if len(p.Images) > MaxProductImages {
		return nil, ErrTooManyProductImages
	}
	if p.Rating == 0 {
		p.Rating = DefaultRating
	}
	if p.Rating < 1 || p.Rating > 5 {
		return nil, ErrInvalidRating
	}
	if p.Stock < 0 {
		return nil, ErrNegativeStock
	}
	if p.DiscountPrice != nil {
		if p.DiscountPrice.Currency() != p.Price.Currency() || p.DiscountPrice.Amount() >= p.Price.Amount() {
			return nil, ErrDiscountTooHigh
		}
	}
	if p.Images == nil {
		p.Images = []ProductImage{}
	}

	now := time.Now()
	return &Product{
		id:             p.ID,
		nameUz:         p.NameUz,
		nameEng:        p.NameEng,
		nameRu:         p.NameRu,
		descriptionUz:  p.DescriptionUz,
		descriptionEng: p.DescriptionEng,
		descriptionRu:  p.DescriptionRu,
		images:         p.Images,
		categoryID:     p.CategoryID,
		price:          p.Price,
		discountPrice:  p.DiscountPrice,
		slug:           p.Slug,
		isAvailable:    p.IsAvailable,
		rating:         p.Rating,
		stock:          p.Stock,
		soldCount:      0,
		tagUz:          p.TagUz,
		tagEng:         p.TagEng,
		tagRu:          p.TagRu,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// ProductFromRepositoryParams - saqlangan mahsulotni bazadan qayta tiklash uchun.
type ProductFromRepositoryParams struct {
	ID             string
	NameUz         string
	NameEng        string
	NameRu         string
	DescriptionUz  string
	DescriptionEng string
	DescriptionRu  string
	Images         []ProductImage
	CategoryID     string
	Price          Money
	DiscountPrice  *Money
	Slug           string
	IsAvailable    bool
	Rating         float64
	Stock          int
	SoldCount      int
	TagUz          *string
	TagEng         *string
	TagRu          *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func NewProductFromRepository(p ProductFromRepositoryParams) *Product {
	return &Product{
		id:             p.ID,
		nameUz:         p.NameUz,
		nameEng:        p.NameEng,
		nameRu:         p.NameRu,
		descriptionUz:  p.DescriptionUz,
		descriptionEng: p.DescriptionEng,
		descriptionRu:  p.DescriptionRu,
		images:         p.Images,
		categoryID:     p.CategoryID,
		price:          p.Price,
		discountPrice:  p.DiscountPrice,
		slug:           p.Slug,
		isAvailable:    p.IsAvailable,
		rating:         p.Rating,
		stock:          p.Stock,
		soldCount:      p.SoldCount,
		tagUz:          p.TagUz,
		tagEng:         p.TagEng,
		tagRu:          p.TagRu,
		createdAt:      p.CreatedAt,
		updatedAt:      p.UpdatedAt,
		deletedAt:      p.DeletedAt,
	}
}

// GenerateSlug - mahsulot nomidan URL uchun qulay slug hosil qiladi.
func GenerateSlug(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	prevDash := true // boshida "-" chiqmasligi uchun
	for _, r := range lower {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// change methods

func (p *Product) ChangeNames(nameUz, nameEng, nameRu string) error {
	if nameUz == "" || nameEng == "" || nameRu == "" {
		return ErrEmptyProductName
	}
	p.nameUz = nameUz
	p.nameEng = nameEng
	p.nameRu = nameRu
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeDescriptions(descriptionUz, descriptionEng, descriptionRu string) {
	p.descriptionUz = descriptionUz
	p.descriptionEng = descriptionEng
	p.descriptionRu = descriptionRu
	p.updatedAt = time.Now()
}

func (p *Product) ChangeImages(images []ProductImage) error {
	if len(images) > MaxProductImages {
		return ErrTooManyProductImages
	}
	p.images = images
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeCategory(categoryID string) error {
	if categoryID == "" {
		return ErrEmptyCategoryID
	}
	p.categoryID = categoryID
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangePrice(price Money) {
	p.price = price
	p.updatedAt = time.Now()
}

func (p *Product) ChangeDiscountPrice(discount *Money) error {
	if discount != nil {
		if discount.Currency() != p.price.Currency() || discount.Amount() >= p.price.Amount() {
			return ErrDiscountTooHigh
		}
	}
	p.discountPrice = discount
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeSlug(slug string) error {
	if slug == "" {
		return ErrEmptyProductSlug
	}
	p.slug = slug
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeAvailability(isAvailable bool) {
	p.isAvailable = isAvailable
	p.updatedAt = time.Now()
}

func (p *Product) ChangeRating(rating float64) error {
	if rating < 1 || rating > 5 {
		return ErrInvalidRating
	}
	p.rating = rating
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeStock(stock int) error {
	if stock < 0 {
		return ErrNegativeStock
	}
	p.stock = stock
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeSoldCount(count int) error {
	if count < 0 {
		return ErrNegativeSoldCount
	}
	p.soldCount = count
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeTags(tagUz, tagEng, tagRu *string) {
	p.tagUz = tagUz
	p.tagEng = tagEng
	p.tagRu = tagRu
	p.updatedAt = time.Now()
}

func (p *Product) Delete() {
	now := time.Now()
	p.deletedAt = &now
}

func (p *Product) IsDeleted() bool {
	return p.deletedAt != nil
}

// getters

func (p *Product) ID() string             { return p.id }
func (p *Product) NameUz() string         { return p.nameUz }
func (p *Product) NameEng() string        { return p.nameEng }
func (p *Product) NameRu() string         { return p.nameRu }
func (p *Product) DescriptionUz() string  { return p.descriptionUz }
func (p *Product) DescriptionEng() string { return p.descriptionEng }
func (p *Product) DescriptionRu() string  { return p.descriptionRu }
func (p *Product) Images() []ProductImage { return p.images }
func (p *Product) CategoryID() string     { return p.categoryID }
func (p *Product) Price() Money           { return p.price }
func (p *Product) DiscountPrice() *Money  { return p.discountPrice }
func (p *Product) Slug() string           { return p.slug }
func (p *Product) IsAvailable() bool      { return p.isAvailable }
func (p *Product) Rating() float64        { return p.rating }
func (p *Product) Stock() int             { return p.stock }
func (p *Product) SoldCount() int         { return p.soldCount }
func (p *Product) TagUz() *string         { return p.tagUz }
func (p *Product) TagEng() *string        { return p.tagEng }
func (p *Product) TagRu() *string         { return p.tagRu }
func (p *Product) CreatedAt() time.Time   { return p.createdAt }
func (p *Product) UpdatedAt() time.Time   { return p.updatedAt }
func (p *Product) DeletedAt() *time.Time  { return p.deletedAt }
