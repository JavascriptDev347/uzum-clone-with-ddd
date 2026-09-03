package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetCategoryUseCase struct {
	repo domain.CategoryRepository
}

func NewGetCategoryUseCase(repo domain.CategoryRepository) *GetCategoryUseCase {
	return &GetCategoryUseCase{repo: repo}
}

func (uc *GetCategoryUseCase) Execute(ctx context.Context, id string, lang Lang) (CategoryOutput, error) {
	category, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return CategoryOutput{}, err
	}

	return *ToCategoryOutput(category, lang), nil
}
