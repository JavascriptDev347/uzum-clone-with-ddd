package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetAllCategoriesIncludingDeletedUseCase struct {
	repo domain.CategoryRepository
}

func NewGetAllCategoriesIncludingDeletedUseCase(repo domain.CategoryRepository) *GetAllCategoriesIncludingDeletedUseCase {
	return &GetAllCategoriesIncludingDeletedUseCase{repo: repo}
}

func (uc *GetAllCategoriesIncludingDeletedUseCase) Execute(ctx context.Context, search string) ([]*CategoryOutputForAdmin, error) {
	categories, err := uc.repo.FindAllIncludingDeleted(ctx, search)
	if err != nil {
		return nil, err
	}

	return ToCategoryOutputsForAdmin(categories), nil
}
