package http

import "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"

type CreateCategoryRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
}

type CreateCategoryResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ParentID  *string `json:"parent_id,omitempty"`
	CreatedAt string  `json:"created_at"`
}

func ToCategoryResponse(output application.CreateCategoryOutput) application.CreateCategoryOutput {
	return application.CreateCategoryOutput{
		ID:        output.ID,
		Name:      output.Name,
		ParentID:  output.ParentID,
		CreatedAt: output.CreatedAt,
	}
}
