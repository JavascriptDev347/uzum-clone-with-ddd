package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/google/uuid"
)

type CreateProductUseCase struct {
	repo domain.ProductRepository
}

func NewCreateProductUseCase(repo domain.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, product CreateProductInput) (CreateProductOutput, error) {
	var id string = uuid.New().String()
	newMoney, err := domain.NewMoney(product.Amount, product.Currency)
	if err != nil {
		return CreateProductOutput{}, err
	}
	newProduct, err := domain.NewProduct(id, product.Name, product.CategoryID, newMoney)
	if err != nil {
		return CreateProductOutput{}, err
	}

	err = uc.repo.Save(ctx, newProduct)
	if err != nil {
		return CreateProductOutput{}, err
	}

	return CreateProductOutput{
		ID:            newProduct.ID(),
		Name:          newProduct.Name(),
		PriceAmount:   newProduct.Price().Amount(),
		PriceCurrency: newProduct.Price().Currency(),
		CategoryID:    newProduct.CategoryID(),
		CreatedAt:     newProduct.CreatedAt(),
	}, nil
}
