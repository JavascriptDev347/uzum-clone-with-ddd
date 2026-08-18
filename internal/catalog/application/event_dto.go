package application

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type CreateEventInput struct {
	Eyebrow    string
	Title      string
	Subtitle   string
	CTA        string
	CategoryID string
	IsRoot     bool
	Image      media.UploadInput
}

type UpdateEventInput struct {
	ID         string  `json:"-"`
	Eyebrow    *string `json:"eyebrow,omitempty"`
	Title      *string `json:"title,omitempty"`
	Subtitle   *string `json:"subtitle,omitempty"`
	CTA        *string `json:"cta,omitempty"`
	CategoryID *string `json:"category_id,omitempty"`
	IsRoot     *bool   `json:"is_root,omitempty"`
}

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

type EventOutputForAdmin struct {
	EventOutput
	DeletedAt *time.Time `json:"deleted_at"`
}

func ToEventOutput(e *domain.Event) *EventOutput {
	return &EventOutput{
		ID:         e.ID(),
		Eyebrow:    e.Eyebrow(),
		Title:      e.Title(),
		Subtitle:   e.Subtitle(),
		CTA:        e.CTA(),
		Image:      e.ImageURL(),
		CategoryID: e.CategoryID(),
		IsRoot:     e.IsRoot(),
		CreatedAt:  e.CreatedAt(),
		UpdatedAt:  e.UpdatedAt(),
	}
}

func ToEventOutputs(events []*domain.Event) []*EventOutput {
	outputs := make([]*EventOutput, 0, len(events))
	for _, e := range events {
		outputs = append(outputs, ToEventOutput(e))
	}
	return outputs
}

func ToEventOutputsForAdmin(events []*domain.Event) []*EventOutputForAdmin {
	outputs := make([]*EventOutputForAdmin, 0, len(events))
	for _, e := range events {
		outputs = append(outputs, &EventOutputForAdmin{
			EventOutput: *ToEventOutput(e),
			DeletedAt:   e.DeletedAt(),
		})
	}
	return outputs
}
