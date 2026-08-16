package http

type CreateProductRequest struct {
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
}

type CreateProductResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	Amount     int64  `json:"amount"`
	ImageURL   string `json:"image_url"`
	Currency   string `json:"currency"`
	CreatedAt  string `json:"created_at"`
}
