package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetProductBySlugUseCase struct {
	repo domain.ProductRepository
}

func NewGetProductBySlugUseCase(repo domain.ProductRepository) *GetProductBySlugUseCase {
	return &GetProductBySlugUseCase{repo: repo}
}

func (uc *GetProductBySlugUseCase) Execute(ctx context.Context, slug string) (*ProductOutput, error) {
	product, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return ToProductOutput(product), nil
}
