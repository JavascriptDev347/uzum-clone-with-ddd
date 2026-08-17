package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type DeleteCategoryUseCase struct {
	repo domain.CategoryRepository
}

func NewDeleteCategoryUseCase(repo domain.CategoryRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{repo: repo}
}

func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, id string) error {
	err := uc.repo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
