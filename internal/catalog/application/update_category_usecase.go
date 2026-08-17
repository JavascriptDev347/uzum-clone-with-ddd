package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type UpdateCategoryUseCase struct {
	repo domain.CategoryRepository
}

func NewUpdateCategoryUseCase(repo domain.CategoryRepository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{repo: repo}
}

func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, input UpdateCategoryInput) error {
	category, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if input.Name != nil {
		if err := category.ChangeName(*input.Name); err != nil {
			return err
		}
	}

	return uc.repo.Update(ctx, category)
}
