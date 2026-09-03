package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/google/uuid"
)

type CreateProductUseCase struct {
	repo         domain.ProductRepository
	categoryRepo domain.CategoryRepository
	uploader     media.Uploader
}

func NewCreateProductUseCase(repo domain.ProductRepository, categoryRepo domain.CategoryRepository, uploader media.Uploader) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo, categoryRepo: categoryRepo, uploader: uploader}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (*ProductOutput, error) {
	if _, err := uc.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	if len(input.Images) > domain.MaxProductImages {
		return nil, domain.ErrTooManyProductImages
	}

	price, err := domain.NewMoney(input.Amount, input.Currency)
	if err != nil {
		return nil, err
	}

	var discountPrice *domain.Money
	if input.DiscountAmount != nil {
		dp, err := domain.NewMoney(*input.DiscountAmount, input.Currency)
		if err != nil {
			return nil, err
		}
		discountPrice = &dp
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

	slug := input.Slug
	if slug == "" {
		slug = domain.GenerateSlug(input.NameUz)
	}
	if slug == "" {
		slug = uuid.New().String()
	}

	id := uuid.New().String()
	product, err := domain.NewProduct(domain.NewProductParams{
		ID:             id,
		NameUz:         input.NameUz,
		NameEng:        input.NameEng,
		NameRu:         input.NameRu,
		DescriptionUz:  input.DescriptionUz,
		DescriptionEng: input.DescriptionEng,
		DescriptionRu:  input.DescriptionRu,
		Images:         uploaded,
		CategoryID:     input.CategoryID,
		Price:          price,
		DiscountPrice:  discountPrice,
		Slug:           slug,
		IsAvailable:    input.IsAvailable,
		Rating:         input.Rating,
		Stock:          input.Stock,
		TagUz:          input.TagUz,
		TagEng:         input.TagEng,
		TagRu:          input.TagRu,
	})
	if err != nil {
		rollback()
		return nil, err
	}

	if err := uc.repo.Save(ctx, product); err != nil {
		rollback()
		return nil, err
	}

	return ToProductOutput(product, LangUZ), nil
}
