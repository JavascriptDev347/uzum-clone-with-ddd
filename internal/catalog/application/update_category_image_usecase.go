package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type UpdateCategoryImageUseCase struct {
	repo     domain.CategoryRepository
	uploader media.Uploader
}

func NewUpdateCategoryImageUseCase(repo domain.CategoryRepository, uploader media.Uploader) *UpdateCategoryImageUseCase {
	return &UpdateCategoryImageUseCase{repo: repo, uploader: uploader}
}

// Execute - kategoriyaning rasmini yangisiga almashtiradi va eski rasmni uploader'dan o'chiradi.
func (uc *UpdateCategoryImageUseCase) Execute(ctx context.Context, input UpdateCategoryImageInput) (*CategoryOutput, error) {
	category, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := media.Validate(input.Image, media.DefaultImageRules()); err != nil {
		return nil, err
	}

	uploadResult, err := uc.uploader.Upload(ctx, input.Image)
	if err != nil {
		return nil, err
	}

	oldPublicID := category.ImagePublicID()
	category.ChangeImage(uploadResult.URL, uploadResult.PublicID)

	if err := uc.repo.Update(ctx, category); err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	if oldPublicID != "" {
		_ = uc.uploader.Delete(ctx, oldPublicID)
	}

	return ToCategoryOutput(category, LangUZ), nil
}
