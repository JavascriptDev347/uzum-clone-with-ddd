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

	if input.NameUz != nil || input.NameEng != nil || input.NameRu != nil {
		nameUz, nameEng, nameRu := category.NameUz(), category.NameEng(), category.NameRu()
		if input.NameUz != nil {
			nameUz = *input.NameUz
		}
		if input.NameEng != nil {
			nameEng = *input.NameEng
		}
		if input.NameRu != nil {
			nameRu = *input.NameRu
		}
		if err := category.ChangeNames(nameUz, nameEng, nameRu); err != nil {
			return err
		}
	}

	return uc.repo.Update(ctx, category)
}
