package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type ReplaceProductImageUseCase struct {
	repo     domain.ProductRepository
	uploader media.Uploader
}

func NewReplaceProductImageUseCase(repo domain.ProductRepository, uploader media.Uploader) *ReplaceProductImageUseCase {
	return &ReplaceProductImageUseCase{repo: repo, uploader: uploader}
}

// Execute - mahsulotning berilgan tartib raqamidagi bitta rasmini yangisiga almashtiradi,
// qolgan rasmlar o'zgarishsiz qoladi. Eski rasm muvaffaqiyatli saqlangandan keyin o'chiriladi.
func (uc *ReplaceProductImageUseCase) Execute(ctx context.Context, input ReplaceProductImageInput) (*ProductOutput, error) {
	product, err := uc.repo.FindByID(ctx, input.ID)
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

	oldImage, err := product.ReplaceImageAt(input.Index, domain.ProductImage{URL: uploadResult.URL, PublicID: uploadResult.PublicID})
	if err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	if err := uc.repo.Update(ctx, product); err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	_ = uc.uploader.Delete(ctx, oldImage.PublicID)

	return ToProductOutput(product, LangUZ), nil
}
