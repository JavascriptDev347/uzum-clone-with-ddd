package application

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

const (
	DefaultProductPage     = 1
	DefaultProductPageSize = 20
	MaxProductPageSize     = 100
)

// NormalizeProductPagination - page/pageSize noto'g'ri yoki bo'sh bo'lsa standart qiymatlarni qo'llaydi.
func NormalizeProductPagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultProductPage
	}
	if pageSize < 1 {
		pageSize = DefaultProductPageSize
	}
	if pageSize > MaxProductPageSize {
		pageSize = MaxProductPageSize
	}
	return page, pageSize
}

// Lang - mahsulot ma'lumotlarini olishda so'raladigan til.
type Lang string

const (
	LangUZ  Lang = "uz"
	LangEng Lang = "eng"
	LangRu  Lang = "ru"
)

// ParseLang - "lang" query parametrini Lang qiymatiga aylantiradi, noto'g'ri/bo'sh bo'lsa uz qaytadi.
func ParseLang(raw string) Lang {
	switch Lang(raw) {
	case LangEng, LangRu:
		return Lang(raw)
	default:
		return LangUZ
	}
}

type CreateProductInput struct {
	NameUz         string
	NameEng        string
	NameRu         string
	DescriptionUz  string
	DescriptionEng string
	DescriptionRu  string
	Images         []media.UploadInput // eng ko'pi bilan 5 ta
	CategoryID     string
	Amount         int64
	Currency       string
	DiscountAmount *int64
	Slug           string
	IsAvailable    bool
	Rating         float64
	Stock          int
	TagUz          *string
	TagEng         *string
	TagRu          *string
}

type UpdateProductInput struct {
	ID             string   `json:"-"`
	NameUz         *string  `json:"name_uz,omitempty"`
	NameEng        *string  `json:"name_eng,omitempty"`
	NameRu         *string  `json:"name_ru,omitempty"`
	DescriptionUz  *string  `json:"description_uz,omitempty"`
	DescriptionEng *string  `json:"description_eng,omitempty"`
	DescriptionRu  *string  `json:"description_ru,omitempty"`
	CategoryID     *string  `json:"category_id,omitempty"`
	Amount         *int64   `json:"amount,omitempty"`
	Currency       *string  `json:"currency,omitempty"`
	DiscountAmount *int64   `json:"discount_amount,omitempty"`
	ClearDiscount  bool     `json:"clear_discount,omitempty"`
	Slug           *string  `json:"slug,omitempty"`
	IsAvailable    *bool    `json:"is_available,omitempty"`
	Rating         *float64 `json:"rating,omitempty"`
	Stock          *int     `json:"stock,omitempty"`
	SoldCount      *int     `json:"sold_count,omitempty"`
	TagUz          *string  `json:"tag_uz,omitempty"`
	ClearTagUz     bool     `json:"clear_tag_uz,omitempty"`
	TagEng         *string  `json:"tag_eng,omitempty"`
	ClearTagEng    bool     `json:"clear_tag_eng,omitempty"`
	TagRu          *string  `json:"tag_ru,omitempty"`
	ClearTagRu     bool     `json:"clear_tag_ru,omitempty"`
}

// ProductOutput - so'ralgan tilga moslashtirilgan (localized) mahsulot ma'lumoti, ommaviy endpointlar uchun.
type ProductOutput struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Tag              *string   `json:"tag,omitempty"`
	Images           []string  `json:"images"`
	CategoryID       string    `json:"category_id"`
	PriceAmount      int64     `json:"price_amount"`
	PriceCurrency    string    `json:"price_currency"`
	DiscountAmount   *int64    `json:"discount_amount,omitempty"`
	FinalPriceAmount int64     `json:"final_price_amount"`
	Slug             string    `json:"slug"`
	IsAvailable      bool      `json:"is_available"`
	Rating           float64   `json:"rating"`
	Stock            int       `json:"stock"`
	SoldCount        int       `json:"sold_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ProductOutputForAdmin - admin panel uchun 3 ta tildagi to'liq ma'lumot (tahrirlash uchun).
type ProductOutputForAdmin struct {
	ID               string     `json:"id"`
	NameUz           string     `json:"name_uz"`
	NameEng          string     `json:"name_eng"`
	NameRu           string     `json:"name_ru"`
	DescriptionUz    string     `json:"description_uz"`
	DescriptionEng   string     `json:"description_eng"`
	DescriptionRu    string     `json:"description_ru"`
	TagUz            *string    `json:"tag_uz,omitempty"`
	TagEng           *string    `json:"tag_eng,omitempty"`
	TagRu            *string    `json:"tag_ru,omitempty"`
	Images           []string   `json:"images"`
	CategoryID       string     `json:"category_id"`
	PriceAmount      int64      `json:"price_amount"`
	PriceCurrency    string     `json:"price_currency"`
	DiscountAmount   *int64     `json:"discount_amount,omitempty"`
	FinalPriceAmount int64      `json:"final_price_amount"`
	Slug             string     `json:"slug"`
	IsAvailable      bool       `json:"is_available"`
	Rating           float64    `json:"rating"`
	Stock            int        `json:"stock"`
	SoldCount        int        `json:"sold_count"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

// localizedField - so'ralgan tildagi qiymat bo'sh bo'lsa uz variantiga fallback qiladi.
func localizedField(uz, eng, ru string, lang Lang) string {
	switch lang {
	case LangEng:
		if eng != "" {
			return eng
		}
	case LangRu:
		if ru != "" {
			return ru
		}
	}
	return uz
}

func localizedTag(uz, eng, ru *string, lang Lang) *string {
	switch lang {
	case LangEng:
		if eng != nil && *eng != "" {
			return eng
		}
	case LangRu:
		if ru != nil && *ru != "" {
			return ru
		}
	}
	return uz
}

func ToProductOutput(p *domain.Product, lang Lang) *ProductOutput {
	images := make([]string, 0, len(p.Images()))
	for _, img := range p.Images() {
		images = append(images, img.URL)
	}

	var discountAmount *int64
	finalPrice := p.Price().Amount()
	if dp := p.DiscountPrice(); dp != nil {
		amt := dp.Amount()
		discountAmount = &amt
		finalPrice = amt
	}

	return &ProductOutput{
		ID:               p.ID(),
		Name:             localizedField(p.NameUz(), p.NameEng(), p.NameRu(), lang),
		Description:      localizedField(p.DescriptionUz(), p.DescriptionEng(), p.DescriptionRu(), lang),
		Tag:              localizedTag(p.TagUz(), p.TagEng(), p.TagRu(), lang),
		Images:           images,
		CategoryID:       p.CategoryID(),
		PriceAmount:      p.Price().Amount(),
		PriceCurrency:    p.Price().Currency(),
		DiscountAmount:   discountAmount,
		FinalPriceAmount: finalPrice,
		Slug:             p.Slug(),
		IsAvailable:      p.IsAvailable(),
		Rating:           p.Rating(),
		Stock:            p.Stock(),
		SoldCount:        p.SoldCount(),
		CreatedAt:        p.CreatedAt(),
		UpdatedAt:        p.UpdatedAt(),
	}
}

func ToProductOutputs(products []*domain.Product, lang Lang) []*ProductOutput {
	outputs := make([]*ProductOutput, 0, len(products))
	for _, p := range products {
		outputs = append(outputs, ToProductOutput(p, lang))
	}
	return outputs
}

func ToProductOutputForAdmin(p *domain.Product) *ProductOutputForAdmin {
	images := make([]string, 0, len(p.Images()))
	for _, img := range p.Images() {
		images = append(images, img.URL)
	}

	var discountAmount *int64
	finalPrice := p.Price().Amount()
	if dp := p.DiscountPrice(); dp != nil {
		amt := dp.Amount()
		discountAmount = &amt
		finalPrice = amt
	}

	return &ProductOutputForAdmin{
		ID:               p.ID(),
		NameUz:           p.NameUz(),
		NameEng:          p.NameEng(),
		NameRu:           p.NameRu(),
		DescriptionUz:    p.DescriptionUz(),
		DescriptionEng:   p.DescriptionEng(),
		DescriptionRu:    p.DescriptionRu(),
		TagUz:            p.TagUz(),
		TagEng:           p.TagEng(),
		TagRu:            p.TagRu(),
		Images:           images,
		CategoryID:       p.CategoryID(),
		PriceAmount:      p.Price().Amount(),
		PriceCurrency:    p.Price().Currency(),
		DiscountAmount:   discountAmount,
		FinalPriceAmount: finalPrice,
		Slug:             p.Slug(),
		IsAvailable:      p.IsAvailable(),
		Rating:           p.Rating(),
		Stock:            p.Stock(),
		SoldCount:        p.SoldCount(),
		CreatedAt:        p.CreatedAt(),
		UpdatedAt:        p.UpdatedAt(),
		DeletedAt:        p.DeletedAt(),
	}
}

func ToProductOutputsForAdmin(products []*domain.Product) []*ProductOutputForAdmin {
	outputs := make([]*ProductOutputForAdmin, 0, len(products))
	for _, p := range products {
		outputs = append(outputs, ToProductOutputForAdmin(p))
	}
	return outputs
}
