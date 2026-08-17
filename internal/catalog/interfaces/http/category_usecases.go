package http

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
)

type CreateCategoryUseCase interface {
	Execute(ctx context.Context, category application.CreateCategoryInput) (*application.CreateCategoryOutput, error)
}

type GetCategoriesUseCase interface {
	Execute(ctx context.Context, search string) ([]*application.CategoryOutput, error)
}

type GetCategoryUseCase interface {
	Execute(ctx context.Context, id string) (application.CategoryOutput, error)
}

type GetCategoryChildUseCase interface {
	Execute(ctx context.Context, parentID string) ([]*application.CategoryOutput, error)
}

type UpdateCategoryUseCase interface {
	Execute(ctx context.Context, input application.UpdateCategoryInput) error
}

type DeleteCategoryUseCase interface {
	Execute(ctx context.Context, id string) error
}
