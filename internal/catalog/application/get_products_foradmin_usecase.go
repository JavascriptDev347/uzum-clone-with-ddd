package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetAllProductsIncludingDeletedUseCase struct {
	repo domain.ProductRepository
}

func NewGetAllProductsIncludingDeletedUseCase(repo domain.ProductRepository) *GetAllProductsIncludingDeletedUseCase {
	return &GetAllProductsIncludingDeletedUseCase{repo: repo}
}

func (uc *GetAllProductsIncludingDeletedUseCase) Execute(ctx context.Context, search string, categoryID string) ([]*ProductOutputForAdmin, error) {
	products, err := uc.repo.FindAllIncludingDeleted(ctx, search, categoryID)
	if err != nil {
		return nil, err
	}
	return ToProductOutputsForAdmin(products), nil
}
