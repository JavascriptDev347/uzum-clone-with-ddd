package application

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type CreateCategoryInput struct {
	Name string `json:"name"`
}

type CreateCategoryOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateCategoryInput struct {
	ID   string
	Name *string
}

type CategoryOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCategoryOutput(c *domain.Category) *CategoryOutput {
	return &CategoryOutput{
		ID:        c.ID(),
		Name:      c.Name(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
}

func ToCategoryOutputs(categories []*domain.Category) []*CategoryOutput {
	outputs := make([]*CategoryOutput, 0, len(categories))
	for _, c := range categories {
		outputs = append(outputs, ToCategoryOutput(c))
	}
	return outputs
}
