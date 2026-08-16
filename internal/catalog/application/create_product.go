package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/google/uuid"
)

type CreateProductUseCase struct {
	repo     domain.ProductRepository
	uploader media.Uploader
}

func NewCreateProductUseCase(repo domain.ProductRepository, uploader media.Uploader) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo, uploader: uploader}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, product CreateProductInput) (CreateProductOutput, error) {
	var id string = uuid.New().String()

	newMoney, err := domain.NewMoney(product.Amount, product.Currency)
	if err != nil {
		return CreateProductOutput{}, err
	}

	var imageURL, imagePublicID string
	if len(product.ImageData) > 0 {
		if err := media.Validate(media.UploadInput{
			ContentType: product.ImageContentType,
			Data:        product.ImageData,
		}, media.DefaultImageRules()); err != nil {
			return CreateProductOutput{}, err
		}

		result, err := uc.uploader.Upload(ctx, media.UploadInput{
			FileName:    product.ImageFileName,
			ContentType: product.ImageContentType,
			Data:        product.ImageData,
		})
		if err != nil {
			return CreateProductOutput{}, err
		}

		imageURL = result.URL
		imagePublicID = result.PublicID
	}

	newProduct, err := domain.NewProduct(id, product.Name, product.CategoryID, newMoney, imageURL, imagePublicID)
	if err != nil {
		return CreateProductOutput{}, err
	}

	if err := uc.repo.Save(ctx, newProduct); err != nil {
		return CreateProductOutput{}, err
	}

	return CreateProductOutput{
		ID:            newProduct.ID(),
		Name:          newProduct.Name(),
		PriceAmount:   newProduct.Price().Amount(),
		PriceCurrency: newProduct.Price().Currency(),
		CategoryID:    newProduct.CategoryID(),
		ImageURL:      newProduct.ImageURL(),
		CreatedAt:     newProduct.CreatedAt(),
	}, nil
}
