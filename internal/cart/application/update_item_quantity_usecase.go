package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type UpdateItemQuantityUseCase struct {
	cartRepo    domain.CartRepository
	productRepo catalog.ProductRepository
}

func NewUpdateItemQuantityUseCase(cartRepo domain.CartRepository, productRepo catalog.ProductRepository) *UpdateItemQuantityUseCase {
	return &UpdateItemQuantityUseCase{cartRepo: cartRepo, productRepo: productRepo}
}

type UpdateItemQuantityInput struct {
	UserID    string
	ProductID string
	Quantity  int // yangi absolyut qiymat — frontend +1/-1 hisoblab shuni yuboradi
}

func (uc *UpdateItemQuantityUseCase) Execute(ctx context.Context, in UpdateItemQuantityInput) error {
	cart, err := uc.cartRepo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return err
	}
	if cart == nil {
		return domain.ErrItemNotFound // cart mavjud emas — item ham bo'lishi mumkin emas
	}

	if in.Quantity <= 0 {
		if err := cart.RemoveItem(in.ProductID); err != nil {
			return err
		}
	} else {
		product, err := uc.productRepo.FindByID(ctx, in.ProductID)
		if err != nil {
			return err
		}
		if product == nil {
			return catalog.ErrProductNotFound
		}
		if in.Quantity > product.Stock() {
			return domain.ErrInsufficientStock
		}

		if err := cart.UpdateItemQuantity(in.ProductID, in.Quantity); err != nil {
			return err
		}
	}

	return uc.cartRepo.Save(ctx, cart)
}
