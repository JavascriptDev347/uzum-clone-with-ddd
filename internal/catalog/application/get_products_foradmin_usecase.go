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

func (uc *GetAllProductsIncludingDeletedUseCase) Execute(ctx context.Context, search string, categoryID string, page, pageSize int) ([]*ProductOutputForAdmin, int64, error) {
	page, pageSize = NormalizeProductPagination(page, pageSize)
	products, total, err := uc.repo.FindAllIncludingDeleted(ctx, search, categoryID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return ToProductOutputsForAdmin(products), total, nil
}
