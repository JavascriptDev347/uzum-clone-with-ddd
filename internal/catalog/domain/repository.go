package domain

import (
	"context"
	"errors"
)

var ErrProductNotFound = errors.New("catalog: product topilmadi")

type ProductRepository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	SoftDelete(ctx context.Context, id string) error
}

type CategoryRepository interface {
	Save(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id string) (*Category, error)
	FindAll(ctx context.Context, search string) ([]*Category, error)
	FindAllIncludingDeleted(ctx context.Context, search string) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	SoftDelete(ctx context.Context, id string) error
}
