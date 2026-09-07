package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type AddProductImagesUseCase struct {
	repo     domain.ProductRepository
	uploader media.Uploader
}

func NewAddProductImagesUseCase(repo domain.ProductRepository, uploader media.Uploader) *AddProductImagesUseCase {
	return &AddProductImagesUseCase{repo: repo, uploader: uploader}
}

// Execute - mahsulotga yangi rasm(lar)ni qo'shadi, mavjud rasmlar o'chirilmaydi.
func (uc *AddProductImagesUseCase) Execute(ctx context.Context, input UpdateProductImagesInput) (*ProductOutput, error) {
	product, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if len(input.Images) == 0 {
		return nil, domain.ErrProductImageRequired
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

	if err := product.AddImages(uploaded); err != nil {
		rollback()
		return nil, err
	}

	if err := uc.repo.Update(ctx, product); err != nil {
		rollback()
		return nil, err
	}

	return ToProductOutput(product, LangUZ), nil
}
