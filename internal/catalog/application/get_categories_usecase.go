package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetCategoriesUseCase struct {
	repo domain.CategoryRepository
}

func NewGetCategoriesUseCase(repo domain.CategoryRepository) *GetCategoriesUseCase {
	return &GetCategoriesUseCase{repo: repo}
}

func (uc *GetCategoriesUseCase) Execute(ctx context.Context, search string) ([]*CategoryOutput, error) {
	categories, err := uc.repo.FindAll(ctx, search)
	if err != nil {
		return nil, err
	}
	return ToCategoryOutputs(categories), nil
}
