package application

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type CreateCategoryInput struct {
	NameUz  string
	NameEng string
	NameRu  string
	Image   media.UploadInput
}

type UpdateCategoryInput struct {
	ID      string
	NameUz  *string `json:"name_uz,omitempty"`
	NameEng *string `json:"name_eng,omitempty"`
	NameRu  *string `json:"name_ru,omitempty"`
}

// CategoryOutput - so'ralgan tilga moslashtirilgan (localized) kategoriya ma'lumoti, ommaviy endpointlar uchun.
type CategoryOutput struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"image_url"`
	ImagePublicID string    `json:"image_public_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateCategoryOutput - kategoriya yaratilgandan keyingi javob, CategoryOutput bilan bir xil shakl.
type CreateCategoryOutput = CategoryOutput

// CategoryOutputForAdmin - admin panel uchun 3 ta tildagi to'liq ma'lumot (tahrirlash uchun).
type CategoryOutputForAdmin struct {
	ID            string     `json:"id"`
	NameUz        string     `json:"name_uz"`
	NameEng       string     `json:"name_eng"`
	NameRu        string     `json:"name_ru"`
	ImageURL      string     `json:"image_url"`
	ImagePublicID string     `json:"image_public_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

func ToCategoryOutput(c *domain.Category, lang Lang) *CategoryOutput {
	return &CategoryOutput{
		ID:            c.ID(),
		Name:          localizedField(c.NameUz(), c.NameEng(), c.NameRu(), lang),
		ImageURL:      c.ImageURL(),
		ImagePublicID: c.ImagePublicID(),
		CreatedAt:     c.CreatedAt(),
		UpdatedAt:     c.UpdatedAt(),
	}
}

func ToCategoryOutputs(categories []*domain.Category, lang Lang) []*CategoryOutput {
	outputs := make([]*CategoryOutput, 0, len(categories))
	for _, c := range categories {
		outputs = append(outputs, ToCategoryOutput(c, lang))
	}
	return outputs
}

func ToCategoryOutputForAdmin(c *domain.Category) *CategoryOutputForAdmin {
	return &CategoryOutputForAdmin{
		ID:            c.ID(),
		NameUz:        c.NameUz(),
		NameEng:       c.NameEng(),
		NameRu:        c.NameRu(),
		ImageURL:      c.ImageURL(),
		ImagePublicID: c.ImagePublicID(),
		CreatedAt:     c.CreatedAt(),
		UpdatedAt:     c.UpdatedAt(),
		DeletedAt:     c.DeletedAt(),
	}
}

func ToCategoryOutputsForAdmin(categories []*domain.Category) []*CategoryOutputForAdmin {
	outputs := make([]*CategoryOutputForAdmin, 0, len(categories))
	for _, c := range categories {
		outputs = append(outputs, ToCategoryOutputForAdmin(c))
	}
	return outputs
}
