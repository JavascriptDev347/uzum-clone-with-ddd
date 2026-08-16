package application

import "time"

type CreateProductInput struct {
	Name             string
	Amount           int64
	Currency         string
	CategoryID       string
	ImageData        []byte
	ImageContentType string
	ImageFileName    string
}

type CreateProductOutput struct {
	ID            string
	Name          string
	PriceAmount   int64
	PriceCurrency string
	CategoryID    string
	ImageURL      string
	CreatedAt     time.Time
}
