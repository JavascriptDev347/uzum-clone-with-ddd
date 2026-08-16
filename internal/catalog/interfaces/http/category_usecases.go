package http

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
)

type CreateCategoryUseCase interface {
	Execute(ctx context.Context, category application.CreateCategoryInput) (application.CreateCategoryOutput, error)
}

type GetCategoriesUseCase interface {
	Execute(ctx context.Context, search string) ([]*application.CategoryOutput, error)
}
