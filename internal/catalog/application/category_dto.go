package application

import "time"

type CreateCategoryInput struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

type CreateCategoryOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
