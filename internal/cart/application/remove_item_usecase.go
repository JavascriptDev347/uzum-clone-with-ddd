package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
)

type RemoveItemUseCase struct {
	cartRepo domain.CartRepository
}

func NewRemoveItemUseCase(cartRepo domain.CartRepository) *RemoveItemUseCase {
	return &RemoveItemUseCase{cartRepo: cartRepo}
}

func (uc *RemoveItemUseCase) Execute(ctx context.Context, userID, productID string) error {
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return domain.ErrItemNotFound
	}

	if err := cart.RemoveItem(productID); err != nil {
		return err
	}

	return uc.cartRepo.Save(ctx, cart)
}
