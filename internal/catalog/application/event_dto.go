package application

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type CreateEventInput struct {
	EyebrowUz  string
	EyebrowEng string
	EyebrowRu  string

	TitleUz  string
	TitleEng string
	TitleRu  string

	SubtitleUz  string
	SubtitleEng string
	SubtitleRu  string

	CTAUz  string
	CTAEng string
	CTARu  string

	CategoryID string
	IsRoot     bool
	Image      media.UploadInput
}

type UpdateEventInput struct {
	ID string `json:"-"`

	EyebrowUz  *string `json:"eyebrow_uz,omitempty"`
	EyebrowEng *string `json:"eyebrow_eng,omitempty"`
	EyebrowRu  *string `json:"eyebrow_ru,omitempty"`

	TitleUz  *string `json:"title_uz,omitempty"`
	TitleEng *string `json:"title_eng,omitempty"`
	TitleRu  *string `json:"title_ru,omitempty"`

	SubtitleUz  *string `json:"subtitle_uz,omitempty"`
	SubtitleEng *string `json:"subtitle_eng,omitempty"`
	SubtitleRu  *string `json:"subtitle_ru,omitempty"`

	CTAUz  *string `json:"cta_uz,omitempty"`
	CTAEng *string `json:"cta_eng,omitempty"`
	CTARu  *string `json:"cta_ru,omitempty"`

	CategoryID *string `json:"category_id,omitempty"`
	IsRoot     *bool   `json:"is_root,omitempty"`
}

// EventOutput - so'ralgan tilga moslashtirilgan (localized) event ma'lumoti, ommaviy endpointlar uchun.
type EventOutput struct {
	ID         string    `json:"id"`
	Eyebrow    string    `json:"eyebrow"`
	Title      string    `json:"title"`
	Subtitle   string    `json:"subtitle"`
	CTA        string    `json:"cta"`
	Image      string    `json:"image"`
	CategoryID string    `json:"category_id"`
	IsRoot     bool      `json:"is_root"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EventOutputForAdmin - admin panel uchun 3 ta tildagi to'liq ma'lumot (tahrirlash uchun).
type EventOutputForAdmin struct {
	ID string `json:"id"`

	EyebrowUz  string `json:"eyebrow_uz"`
	EyebrowEng string `json:"eyebrow_eng"`
	EyebrowRu  string `json:"eyebrow_ru"`

	TitleUz  string `json:"title_uz"`
	TitleEng string `json:"title_eng"`
	TitleRu  string `json:"title_ru"`

	SubtitleUz  string `json:"subtitle_uz"`
	SubtitleEng string `json:"subtitle_eng"`
	SubtitleRu  string `json:"subtitle_ru"`

	CTAUz  string `json:"cta_uz"`
	CTAEng string `json:"cta_eng"`
	CTARu  string `json:"cta_ru"`

	Image      string     `json:"image"`
	CategoryID string     `json:"category_id"`
	IsRoot     bool       `json:"is_root"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func ToEventOutput(e *domain.Event, lang Lang) *EventOutput {
	return &EventOutput{
		ID:         e.ID(),
		Eyebrow:    localizedField(e.EyebrowUz(), e.EyebrowEng(), e.EyebrowRu(), lang),
		Title:      localizedField(e.TitleUz(), e.TitleEng(), e.TitleRu(), lang),
		Subtitle:   localizedField(e.SubtitleUz(), e.SubtitleEng(), e.SubtitleRu(), lang),
		CTA:        localizedField(e.CTAUz(), e.CTAEng(), e.CTARu(), lang),
		Image:      e.ImageURL(),
		CategoryID: e.CategoryID(),
		IsRoot:     e.IsRoot(),
		CreatedAt:  e.CreatedAt(),
		UpdatedAt:  e.UpdatedAt(),
	}
}

func ToEventOutputs(events []*domain.Event, lang Lang) []*EventOutput {
	outputs := make([]*EventOutput, 0, len(events))
	for _, e := range events {
		outputs = append(outputs, ToEventOutput(e, lang))
	}
	return outputs
}

func ToEventOutputForAdmin(e *domain.Event) *EventOutputForAdmin {
	return &EventOutputForAdmin{
		ID:          e.ID(),
		EyebrowUz:   e.EyebrowUz(),
		EyebrowEng:  e.EyebrowEng(),
		EyebrowRu:   e.EyebrowRu(),
		TitleUz:     e.TitleUz(),
		TitleEng:    e.TitleEng(),
		TitleRu:     e.TitleRu(),
		SubtitleUz:  e.SubtitleUz(),
		SubtitleEng: e.SubtitleEng(),
		SubtitleRu:  e.SubtitleRu(),
		CTAUz:       e.CTAUz(),
		CTAEng:      e.CTAEng(),
		CTARu:       e.CTARu(),
		Image:       e.ImageURL(),
		CategoryID:  e.CategoryID(),
		IsRoot:      e.IsRoot(),
		CreatedAt:   e.CreatedAt(),
		UpdatedAt:   e.UpdatedAt(),
		DeletedAt:   e.DeletedAt(),
	}
}

func ToEventOutputsForAdmin(events []*domain.Event) []*EventOutputForAdmin {
	outputs := make([]*EventOutputForAdmin, 0, len(events))
	for _, e := range events {
		outputs = append(outputs, ToEventOutputForAdmin(e))
	}
	return outputs
}
