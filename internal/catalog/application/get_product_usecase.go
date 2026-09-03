package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetProductUseCase struct {
	repo domain.ProductRepository
}

func NewGetProductUseCase(repo domain.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{repo: repo}
}

func (uc *GetProductUseCase) Execute(ctx context.Context, id string, lang Lang) (*ProductOutput, error) {
	product, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToProductOutput(product, lang), nil
}
