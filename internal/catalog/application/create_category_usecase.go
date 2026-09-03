package application

import (
	"context"
	"fmt"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/google/uuid"
)

type CreateCategoryUseCase struct {
	repo     domain.CategoryRepository
	uploader media.Uploader
}

func NewCreateCategoryUseCase(repo domain.CategoryRepository, uploader media.Uploader) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo, uploader: uploader}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) (*CreateCategoryOutput, error) {
	uploadResult, err := uc.uploader.Upload(ctx, input.Image)
	if err != nil {
		return nil, fmt.Errorf("category image upload failed: %w", err)
	}

	id := uuid.New().String()

	category, err := domain.NewCategory(id, input.NameUz, input.NameEng, input.NameRu, uploadResult.URL, uploadResult.PublicID)
	if err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	if err := uc.repo.Save(ctx, category); err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	return ToCategoryOutput(category, LangUZ), nil
}
