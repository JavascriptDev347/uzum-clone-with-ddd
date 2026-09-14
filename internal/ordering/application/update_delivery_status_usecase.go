package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// UpdateDeliveryStatusUseCase - admin uchun: buyurtmaning yetkazib berish holatini
// o'zgartiradi. Orqaga qaytarishni taqiqlash mantig'i bu yerda takrorlanmaydi - u to'liq
// Order.AdvanceDelivery()ga yuklangan (domain aggregate o'z invariantini o'zi himoya qiladi).
type UpdateDeliveryStatusUseCase struct {
	orderRepo domain.OrderRepository
}

func NewUpdateDeliveryStatusUseCase(orderRepo domain.OrderRepository) *UpdateDeliveryStatusUseCase {
	return &UpdateDeliveryStatusUseCase{orderRepo: orderRepo}
}

func (uc *UpdateDeliveryStatusUseCase) Execute(ctx context.Context, orderID string, status domain.DeliveryStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	if err := order.AdvanceDelivery(status); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
