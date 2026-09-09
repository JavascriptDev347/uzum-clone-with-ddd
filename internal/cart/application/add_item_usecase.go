package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type AddItemUseCase struct {
	cartRepo    domain.CartRepository
	productRepo catalog.ProductRepository // TODO: haqiqiy interfeys nomini tekshir
}

func NewAddItemUseCase(cartRepo domain.CartRepository, productRepo catalog.ProductRepository) *AddItemUseCase {
	return &AddItemUseCase{cartRepo: cartRepo, productRepo: productRepo}
}

type AddItemInput struct {
	UserID    string
	ProductID string
	Quantity  int
}

func (uc *AddItemUseCase) Execute(ctx context.Context, in AddItemInput) error {
	// TODO: haqiqiy metod nomini tekshir — GetByID deb taxmin qildim
	product, err := uc.productRepo.FindByID(ctx, in.ProductID)
	if err != nil {
		return err
	}
	if product == nil {
		return catalog.ErrProductNotFound // TODO: haqiqiy xato nomini tekshir
	}

	cart, err := uc.cartRepo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return err
	}
	if cart == nil {
		cart = domain.NewCart(in.UserID)
	}

	if cart.QuantityOf(in.ProductID)+in.Quantity > product.Stock() {
		return domain.ErrInsufficientStock
	}

	if err := cart.AddItem(in.ProductID, in.Quantity); err != nil {
		return err
	}

	return uc.cartRepo.Save(ctx, cart)
}
