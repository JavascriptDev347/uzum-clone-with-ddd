package domain

import (
	"context"
	"errors"
)

var ErrProductNotFound = errors.New("catalog: product topilmadi")

type ProductRepository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindBySlug(ctx context.Context, slug string) (*Product, error)
	FindAll(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*Product, int64, error)
	FindAllIncludingDeleted(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*Product, int64, error)
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

type EventRepository interface {
	Save(ctx context.Context, event *Event) error
	FindByID(ctx context.Context, id string) (*Event, error)
	FindAll(ctx context.Context) ([]*Event, error)
	FindAllIncludingDeleted(ctx context.Context) ([]*Event, error)
	Update(ctx context.Context, event *Event) error
	SoftDelete(ctx context.Context, id string) error
}
