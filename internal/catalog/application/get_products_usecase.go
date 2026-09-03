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

func (uc *GetProductsUseCase) Execute(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*ProductOutput, int64, error) {
	page, pageSize = NormalizeProductPagination(page, pageSize)
	products, total, err := uc.repo.FindAll(ctx, search, categoryID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return ToProductOutputs(products), total, nil
}
