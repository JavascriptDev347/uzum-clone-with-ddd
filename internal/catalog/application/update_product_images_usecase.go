package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type UpdateProductImagesUseCase struct {
	repo     domain.ProductRepository
	uploader media.Uploader
}

func NewUpdateProductImagesUseCase(repo domain.ProductRepository, uploader media.Uploader) *UpdateProductImagesUseCase {
	return &UpdateProductImagesUseCase{repo: repo, uploader: uploader}
}

// Execute - mahsulotning rasmlar to'plamini butunlay yangisiga almashtiradi va eski rasmlarni uploader'dan o'chiradi.
func (uc *UpdateProductImagesUseCase) Execute(ctx context.Context, input UpdateProductImagesInput) (*ProductOutput, error) {
	product, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if len(input.Images) == 0 {
		return nil, domain.ErrProductImageRequired
	}
	if len(input.Images) > domain.MaxProductImages {
		return nil, domain.ErrTooManyProductImages
	}

	uploaded := make([]domain.ProductImage, 0, len(input.Images))
	rollback := func() {
		for _, img := range uploaded {
			_ = uc.uploader.Delete(ctx, img.PublicID)
		}
	}

	for _, img := range input.Images {
		if err := media.Validate(img, media.DefaultImageRules()); err != nil {
			rollback()
			return nil, err
		}
		result, err := uc.uploader.Upload(ctx, img)
		if err != nil {
			rollback()
			return nil, err
		}
		uploaded = append(uploaded, domain.ProductImage{URL: result.URL, PublicID: result.PublicID})
	}

	oldImages := product.Images()

	if err := product.ChangeImages(uploaded); err != nil {
		rollback()
		return nil, err
	}

	if err := uc.repo.Update(ctx, product); err != nil {
		rollback()
		return nil, err
	}

	for _, img := range oldImages {
		_ = uc.uploader.Delete(ctx, img.PublicID)
	}

	return ToProductOutput(product, LangUZ), nil
}
