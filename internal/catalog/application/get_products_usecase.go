package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetProductsUseCase struct {
	repo domain.ProductRepository
}

func NewGetProductsUseCase(repo domain.ProductRepository) *GetProductsUseCase {
	return &GetProductsUseCase{repo: repo}
}

func (uc *GetProductsUseCase) Execute(ctx context.Context, search string, categoryID string) ([]*ProductOutput, error) {
	products, err := uc.repo.FindAll(ctx, search, categoryID)
	if err != nil {
		return nil, err
	}
	return ToProductOutputs(products), nil
}
