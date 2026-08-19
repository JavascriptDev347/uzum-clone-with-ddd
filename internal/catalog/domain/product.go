package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyProductID       = errors.New("catalog: product ID bo'sh bo'lishi mumkin emas")
	ErrEmptyProductName     = errors.New("catalog: product nomi bo'sh bo'lishi mumkin emas")
	ErrEmptyProductSlug     = errors.New("catalog: product slug bo'sh bo'lishi mumkin emas")
	ErrEmptyCategoryID      = errors.New("catalog: category ID bo'sh bo'lishi mumkin emas")
	ErrTooManyProductImages = errors.New("catalog: mahsulot uchun eng ko'pi bilan 5 ta rasm yuklash mumkin")
	ErrInvalidRating        = errors.New("catalog: reyting 1 dan 5 gacha bo'lishi kerak")
	ErrNegativeStock        = errors.New("catalog: stock manfiy bo'lishi mumkin emas")
	ErrNegativeSoldCount    = errors.New("catalog: sotilgan mahsulotlar soni manfiy bo'lishi mumkin emas")
	ErrNegativeStemCount    = errors.New("catalog: gullar soni manfiy bo'lishi mumkin emas")
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
	id                string
	name              string
	description       string
	images            []ProductImage
	videoURLYoutube   string // YouTube video URL, frontendda embed qilinadi
	videoURLInstagram string // Instagram video URL, frontendda embed qilinadi
	categoryID        string

	price         Money
	discountPrice *Money // ixtiyoriy, bo'lsa umumiy summa shu bo'yicha hisoblanadi

	slug        string
	isAvailable bool
	rating      float64 // 1-5, default 1

	stock     int
	soldCount int

	flowerTypes []string // gullar turi (tag sifatida)
	color       string
	stemCount   int

	packagingType     PackagingType     // enum: bucket, box, vase
	freshnessLifespan FreshnessLifespan // kunlarda, 1-7, default 1

	careInstructions *string  // ixtiyoriy, default nil
	occasions        []string // taglar: masalan tug'ilgan kun, to'y va h.k.

	allowCustomCard  *bool    // hozircha har doim nil - keyinchalik qo'shiladi
	compatibleAddons []string // ixtiyoriy

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

// NewProductParams - yangi mahsulot yaratish uchun kerakli ma'lumotlar.
type NewProductParams struct {
	ID                string
	Name              string
	Description       string
	Images            []ProductImage
	VideoURLYoutube   string
	VideoURLInstagram string
	CategoryID        string
	Price             Money
	DiscountPrice     *Money
	Slug              string
	IsAvailable       bool
	Rating            float64
	Stock             int
	FlowerTypes       []string
	Color             string
	StemCount         int
	PackagingType     PackagingType
	FreshnessLifespan FreshnessLifespan
	CareInstructions  *string
	Occasions         []string
	CompatibleAddons  []string
}

func NewProduct(p NewProductParams) (*Product, error) {
	if p.ID == "" {
		return nil, ErrEmptyProductID
	}
	if p.Name == "" {
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
	if !p.PackagingType.IsValid() {
		return nil, ErrInvalidPackagingType
	}
	if p.FreshnessLifespan == 0 {
		p.FreshnessLifespan = DefaultFreshnessLifespan
	}
	if !p.FreshnessLifespan.IsValid() {
		return nil, ErrInvalidFreshnessLifespan
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
	if p.StemCount < 0 {
		return nil, ErrNegativeStemCount
	}
	if p.DiscountPrice != nil {
		if p.DiscountPrice.Currency() != p.Price.Currency() || p.DiscountPrice.Amount() >= p.Price.Amount() {
			return nil, ErrDiscountTooHigh
		}
	}
	if p.FlowerTypes == nil {
		p.FlowerTypes = []string{}
	}
	if p.Occasions == nil {
		p.Occasions = []string{}
	}
	if p.CompatibleAddons == nil {
		p.CompatibleAddons = []string{}
	}
	if p.Images == nil {
		p.Images = []ProductImage{}
	}

	now := time.Now()
	return &Product{
		id:                p.ID,
		name:              p.Name,
		description:       p.Description,
		images:            p.Images,
		videoURLYoutube:   p.VideoURLYoutube,
		videoURLInstagram: p.VideoURLInstagram,
		categoryID:        p.CategoryID,
		price:             p.Price,
		discountPrice:     p.DiscountPrice,
		slug:              p.Slug,
		isAvailable:       p.IsAvailable,
		rating:            p.Rating,
		stock:             p.Stock,
		soldCount:         0,
		flowerTypes:       p.FlowerTypes,
		color:             p.Color,
		stemCount:         p.StemCount,
		packagingType:     p.PackagingType,
		freshnessLifespan: p.FreshnessLifespan,
		careInstructions:  p.CareInstructions,
		occasions:         p.Occasions,
		allowCustomCard:   nil,
		compatibleAddons:  p.CompatibleAddons,
		createdAt:         now,
		updatedAt:         now,
	}, nil
}

// ProductFromRepositoryParams - saqlangan mahsulotni bazadan qayta tiklash uchun.
type ProductFromRepositoryParams struct {
	ID                string
	Name              string
	Description       string
	Images            []ProductImage
	VideoURLYoutube   string
	VideoURLInstagram string
	CategoryID        string
	Price             Money
	DiscountPrice     *Money
	Slug              string
	IsAvailable       bool
	Rating            float64
	Stock             int
	SoldCount         int
	FlowerTypes       []string
	Color             string
	StemCount         int
	PackagingType     PackagingType
	FreshnessLifespan FreshnessLifespan
	CareInstructions  *string
	Occasions         []string
	AllowCustomCard   *bool
	CompatibleAddons  []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

func NewProductFromRepository(p ProductFromRepositoryParams) *Product {
	return &Product{
		id:                p.ID,
		name:              p.Name,
		description:       p.Description,
		images:            p.Images,
		videoURLYoutube:   p.VideoURLYoutube,
		videoURLInstagram: p.VideoURLInstagram,
		categoryID:        p.CategoryID,
		price:             p.Price,
		discountPrice:     p.DiscountPrice,
		slug:              p.Slug,
		isAvailable:       p.IsAvailable,
		rating:            p.Rating,
		stock:             p.Stock,
		soldCount:         p.SoldCount,
		flowerTypes:       p.FlowerTypes,
		color:             p.Color,
		stemCount:         p.StemCount,
		packagingType:     p.PackagingType,
		freshnessLifespan: p.FreshnessLifespan,
		careInstructions:  p.CareInstructions,
		occasions:         p.Occasions,
		allowCustomCard:   p.AllowCustomCard,
		compatibleAddons:  p.CompatibleAddons,
		createdAt:         p.CreatedAt,
		updatedAt:         p.UpdatedAt,
		deletedAt:         p.DeletedAt,
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

func (p *Product) ChangeName(name string) error {
	if name == "" {
		return ErrEmptyProductName
	}
	p.name = name
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeDescription(description string) {
	p.description = description
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

func (p *Product) ChangeVideoURLYoutube(url string) {
	p.videoURLYoutube = url
	p.updatedAt = time.Now()
}

func (p *Product) ChangeVideoURLInstagram(url string) {
	p.videoURLInstagram = url
	p.updatedAt = time.Now()
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

func (p *Product) ChangeFlowerTypes(types []string) {
	if types == nil {
		types = []string{}
	}
	p.flowerTypes = types
	p.updatedAt = time.Now()
}

func (p *Product) ChangeColor(color string) {
	p.color = color
	p.updatedAt = time.Now()
}

func (p *Product) ChangeStemCount(count int) error {
	if count < 0 {
		return ErrNegativeStemCount
	}
	p.stemCount = count
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangePackagingType(pt PackagingType) error {
	if !pt.IsValid() {
		return ErrInvalidPackagingType
	}
	p.packagingType = pt
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeFreshnessLifespan(f FreshnessLifespan) error {
	if !f.IsValid() {
		return ErrInvalidFreshnessLifespan
	}
	p.freshnessLifespan = f
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeCareInstructions(instructions *string) {
	p.careInstructions = instructions
	p.updatedAt = time.Now()
}

func (p *Product) ChangeOccasions(occasions []string) {
	if occasions == nil {
		occasions = []string{}
	}
	p.occasions = occasions
	p.updatedAt = time.Now()
}

func (p *Product) ChangeCompatibleAddons(addons []string) {
	if addons == nil {
		addons = []string{}
	}
	p.compatibleAddons = addons
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

func (p *Product) ID() string                           { return p.id }
func (p *Product) Name() string                         { return p.name }
func (p *Product) Description() string                  { return p.description }
func (p *Product) Images() []ProductImage               { return p.images }
func (p *Product) VideoURLYoutube() string              { return p.videoURLYoutube }
func (p *Product) VideoURLInstagram() string            { return p.videoURLInstagram }
func (p *Product) CategoryID() string                   { return p.categoryID }
func (p *Product) Price() Money                         { return p.price }
func (p *Product) DiscountPrice() *Money                { return p.discountPrice }
func (p *Product) Slug() string                         { return p.slug }
func (p *Product) IsAvailable() bool                    { return p.isAvailable }
func (p *Product) Rating() float64                      { return p.rating }
func (p *Product) Stock() int                           { return p.stock }
func (p *Product) SoldCount() int                       { return p.soldCount }
func (p *Product) FlowerTypes() []string                { return p.flowerTypes }
func (p *Product) Color() string                        { return p.color }
func (p *Product) StemCount() int                       { return p.stemCount }
func (p *Product) PackagingType() PackagingType         { return p.packagingType }
func (p *Product) FreshnessLifespan() FreshnessLifespan { return p.freshnessLifespan }
func (p *Product) CareInstructions() *string            { return p.careInstructions }
func (p *Product) Occasions() []string                  { return p.occasions }
func (p *Product) AllowCustomCard() *bool               { return p.allowCustomCard }
func (p *Product) CompatibleAddons() []string           { return p.compatibleAddons }
func (p *Product) CreatedAt() time.Time                 { return p.createdAt }
func (p *Product) UpdatedAt() time.Time                 { return p.updatedAt }
func (p *Product) DeletedAt() *time.Time                { return p.deletedAt }
