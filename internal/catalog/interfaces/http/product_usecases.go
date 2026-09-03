package http

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
)

type CreateProductUseCase interface {
	Execute(ctx context.Context, input application.CreateProductInput) (*application.ProductOutput, error)
}

type GetProductsUseCase interface {
	Execute(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*application.ProductOutput, int64, error)
}

type GetProductUseCase interface {
	Execute(ctx context.Context, id string) (*application.ProductOutput, error)
}

type GetProductBySlugUseCase interface {
	Execute(ctx context.Context, slug string) (*application.ProductOutput, error)
}

type UpdateProductUseCase interface {
	Execute(ctx context.Context, input application.UpdateProductInput) error
}

type DeleteProductUseCase interface {
	Execute(ctx context.Context, id string) error
}

type GetAllProductsIncludingDeletedUseCase interface {
	Execute(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*application.ProductOutputForAdmin, int64, error)
}
