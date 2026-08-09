package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyProductID   = errors.New("catalog: product ID bo'sh bo'lishi mumkin emas")
	ErrEmptyProductName = errors.New("catalog: product nomi bo'sh bo'lishi mumkin emas")
	ErrEmptyCategoryID  = errors.New("catalog: category ID bo'sh bo'lishi mumkin emas")
)

type Product struct {
	id         string
	name       string
	categoryID string
	price      Money
	createdAt  time.Time
}

func NewProduct(id, name, categoryID string, price Money) (*Product, error) {
	if id == "" {
		return nil, ErrEmptyProductID
	}
	if name == "" {
		return nil, ErrEmptyProductName
	}
	if categoryID == "" {
		return nil, ErrEmptyCategoryID
	}

	return &Product{
		id:         id,
		name:       name,
		categoryID: categoryID,
		price:      price,
		createdAt:  time.Now(),
	}, nil
}

func NewProductFromRepository(id, name, categoryID string, price Money, createdAt time.Time) *Product {
	return &Product{
		id:         id,
		name:       name,
		categoryID: categoryID,
		price:      price,
		createdAt:  createdAt,
	}
}

func (p *Product) ChangePrice(price Money) {
	p.price = price
}

func (p *Product) ChangeName(name string) error {
	if name == "" {
		return ErrEmptyProductName
	}
	p.name = name
	return nil
}

// getters
func (p *Product) ID() string {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) CategoryID() string {
	return p.categoryID
}

func (p *Product) Price() Money {
	return p.price
}

func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}
