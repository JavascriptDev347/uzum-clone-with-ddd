package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// UpdatePaymentStatusUseCase - admin uchun: buyurtmaning to'lov holatini o'zgartiradi.
type UpdatePaymentStatusUseCase struct {
	orderRepo domain.OrderRepository
}

func NewUpdatePaymentStatusUseCase(orderRepo domain.OrderRepository) *UpdatePaymentStatusUseCase {
	return &UpdatePaymentStatusUseCase{orderRepo: orderRepo}
}

func (uc *UpdatePaymentStatusUseCase) Execute(ctx context.Context, orderID string, status domain.PaymentStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	switch status {
	case domain.PaymentStatusPaid:
		order.MarkPaid()
	case domain.PaymentStatusUnpaid:
		order.MarkUnpaid()
	default:
		return nil, domain.ErrInvalidPaymentStatus
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
