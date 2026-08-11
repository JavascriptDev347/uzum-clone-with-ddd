package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/google/uuid"
)

type CreateCategoryUseCase struct {
	repo domain.CategoryRepository
}

func NewCreateCategoryUseCase(repo domain.CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) (CreateCategoryOutput, error) {
	var id string = uuid.New().String()
	newCategory, err := domain.NewCategory(id, input.Name, input.ParentID)
	if err != nil {
		return CreateCategoryOutput{}, err
	}

	err = uc.repo.Save(ctx, newCategory)
	if err != nil {
		return CreateCategoryOutput{}, err
	}

	return CreateCategoryOutput{
		ID:        newCategory.ID(),
		Name:      newCategory.Name(),
		ParentID:  newCategory.ParentID(),
		CreatedAt: newCategory.CreatedAt(),
	}, nil
}
