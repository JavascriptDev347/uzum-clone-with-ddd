package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
)

// ClearCartUseCase - checkout tugagandan keyin cart'ni bo'shatish uchun (ordering context
// buni cart'ning application-layer chegarasi sifatida chaqiradi, cart domain/repository'ga
// to'g'ridan-to'g'ri kirmasdan).
type ClearCartUseCase struct {
	cartRepo domain.CartRepository
}

func NewClearCartUseCase(cartRepo domain.CartRepository) *ClearCartUseCase {
	return &ClearCartUseCase{cartRepo: cartRepo}
}

func (uc *ClearCartUseCase) Execute(ctx context.Context, userID string) error {
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil // cart mavjud emas — tozalanadigan narsa yo'q
	}

	cart.Clear()
	return uc.cartRepo.Save(ctx, cart)
}
