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

type CreateProductInput struct {
	Name              string
	Description       string
	Images            []media.UploadInput // eng ko'pi bilan 5 ta
	VideoURLYoutube   string
	VideoURLInstagram string
	CategoryID        string
	Amount            int64
	Currency          string
	DiscountAmount    *int64
	Slug              string
	IsAvailable       bool
	Rating            float64
	Stock             int
	FlowerTypes       []string
	Color             string
	StemCount         int
	PackagingType     string
	FreshnessLifespan int
	CareInstructions  *string
	Occasions         []string
	CompatibleAddons  []string
}

type UpdateProductInput struct {
	ID                    string   `json:"-"`
	Name                  *string  `json:"name,omitempty"`
	Description           *string  `json:"description,omitempty"`
	VideoURLYoutube       *string  `json:"video_url_youtube,omitempty"`
	VideoURLInstagram     *string  `json:"video_url_instagram,omitempty"`
	CategoryID            *string  `json:"category_id,omitempty"`
	Amount                *int64   `json:"amount,omitempty"`
	Currency              *string  `json:"currency,omitempty"`
	DiscountAmount        *int64   `json:"discount_amount,omitempty"`
	ClearDiscount         bool     `json:"clear_discount,omitempty"`
	Slug                  *string  `json:"slug,omitempty"`
	IsAvailable           *bool    `json:"is_available,omitempty"`
	Rating                *float64 `json:"rating,omitempty"`
	Stock                 *int     `json:"stock,omitempty"`
	SoldCount             *int     `json:"sold_count,omitempty"`
	FlowerTypes           []string `json:"flower_types,omitempty"`
	Color                 *string  `json:"color,omitempty"`
	StemCount             *int     `json:"stem_count,omitempty"`
	PackagingType         *string  `json:"packaging_type,omitempty"`
	FreshnessLifespan     *int     `json:"freshness_lifespan,omitempty"`
	CareInstructions      *string  `json:"care_instructions,omitempty"`
	ClearCareInstructions bool     `json:"clear_care_instructions,omitempty"`
	Occasions             []string `json:"occasions,omitempty"`
	CompatibleAddons      []string `json:"compatible_addons,omitempty"`
}

type ProductOutput struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Images            []string  `json:"images"`
	VideoURLYoutube   string    `json:"video_url_youtube"`
	VideoURLInstagram string    `json:"video_url_instagram"`
	CategoryID        string    `json:"category_id"`
	PriceAmount       int64     `json:"price_amount"`
	PriceCurrency     string    `json:"price_currency"`
	DiscountAmount    *int64    `json:"discount_amount,omitempty"`
	FinalPriceAmount  int64     `json:"final_price_amount"`
	Slug              string    `json:"slug"`
	IsAvailable       bool      `json:"is_available"`
	Rating            float64   `json:"rating"`
	Stock             int       `json:"stock"`
	SoldCount         int       `json:"sold_count"`
	FlowerTypes       []string  `json:"flower_types"`
	Color             string    `json:"color"`
	StemCount         int       `json:"stem_count"`
	PackagingType     string    `json:"packaging_type"`
	FreshnessLifespan int       `json:"freshness_lifespan"`
	CareInstructions  *string   `json:"care_instructions"`
	Occasions         []string  `json:"occasions"`
	AllowCustomCard   *bool     `json:"allow_custom_card"`
	CompatibleAddons  []string  `json:"compatible_addons"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ProductOutputForAdmin struct {
	ProductOutput
	DeletedAt *time.Time `json:"deleted_at"`
}

func ToProductOutput(p *domain.Product) *ProductOutput {
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
		ID:                p.ID(),
		Name:              p.Name(),
		Description:       p.Description(),
		Images:            images,
		VideoURLYoutube:   p.VideoURLYoutube(),
		VideoURLInstagram: p.VideoURLInstagram(),
		CategoryID:        p.CategoryID(),
		PriceAmount:       p.Price().Amount(),
		PriceCurrency:     p.Price().Currency(),
		DiscountAmount:    discountAmount,
		FinalPriceAmount:  finalPrice,
		Slug:              p.Slug(),
		IsAvailable:       p.IsAvailable(),
		Rating:            p.Rating(),
		Stock:             p.Stock(),
		SoldCount:         p.SoldCount(),
		FlowerTypes:       p.FlowerTypes(),
		Color:             p.Color(),
		StemCount:         p.StemCount(),
		PackagingType:     p.PackagingType().String(),
		FreshnessLifespan: int(p.FreshnessLifespan()),
		CareInstructions:  p.CareInstructions(),
		Occasions:         p.Occasions(),
		AllowCustomCard:   p.AllowCustomCard(),
		CompatibleAddons:  p.CompatibleAddons(),
		CreatedAt:         p.CreatedAt(),
		UpdatedAt:         p.UpdatedAt(),
	}
}

func ToProductOutputs(products []*domain.Product) []*ProductOutput {
	outputs := make([]*ProductOutput, 0, len(products))
	for _, p := range products {
		outputs = append(outputs, ToProductOutput(p))
	}
	return outputs
}

func ToProductOutputsForAdmin(products []*domain.Product) []*ProductOutputForAdmin {
	outputs := make([]*ProductOutputForAdmin, 0, len(products))
	for _, p := range products {
		outputs = append(outputs, &ProductOutputForAdmin{
			ProductOutput: *ToProductOutput(p),
			DeletedAt:     p.DeletedAt(),
		})
	}
	return outputs
}
