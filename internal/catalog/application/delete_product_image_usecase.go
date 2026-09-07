package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type DeleteProductImageUseCase struct {
	repo     domain.ProductRepository
	uploader media.Uploader
}

func NewDeleteProductImageUseCase(repo domain.ProductRepository, uploader media.Uploader) *DeleteProductImageUseCase {
	return &DeleteProductImageUseCase{repo: repo, uploader: uploader}
}

// Execute - mahsulotning berilgan tartib raqamidagi bitta rasmini o'chiradi, qolgan rasmlar o'zgarmaydi.
// Rasm muvaffaqiyatli saqlangandan (bazadan) keyin uploader'dan (S3) ham o'chiriladi.
func (uc *DeleteProductImageUseCase) Execute(ctx context.Context, input DeleteProductImageInput) (*ProductOutput, error) {
	product, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	removed, err := product.RemoveImageAt(input.Index)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	_ = uc.uploader.Delete(ctx, removed.PublicID)

	return ToProductOutput(product, LangUZ), nil
}
