package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type DeleteProductUseCase struct {
	repo domain.ProductRepository
}

func NewDeleteProductUseCase(repo domain.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{repo: repo}
}

func (uc *DeleteProductUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
