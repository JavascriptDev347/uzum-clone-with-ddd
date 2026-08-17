package http

import "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CreateCategoryResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ParentID  *string `json:"parent_id,omitempty"`
	UpdatedAt string  `json:"updated_at"`
	CreatedAt string  `json:"created_at"`
}

func ToCategoryResponse(output application.CreateCategoryOutput) application.CreateCategoryOutput {
	return application.CreateCategoryOutput{
		ID:        output.ID,
		Name:      output.Name,
		UpdatedAt: output.UpdatedAt,
		CreatedAt: output.CreatedAt,
	}
}
