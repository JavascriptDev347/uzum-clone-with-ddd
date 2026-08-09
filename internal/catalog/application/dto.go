package application

import "time"

type CreateProductInput struct {
	Name       string `json:"name"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	CategoryID string `json:"category_id"`
}

type CreateProductOutput struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PriceAmount   int64     `json:"price_amount"`
	PriceCurrency string    `json:"price_currency"`
	CategoryID    string    `json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
}
