package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// GetUserOrdersUseCase - foydalanuvchining o'z buyurtmalari ro'yxatini olish uchun
// (buyurtma holati sahifasi uchun).
type GetUserOrdersUseCase struct {
	orderRepo domain.OrderRepository
}

func NewGetUserOrdersUseCase(orderRepo domain.OrderRepository) *GetUserOrdersUseCase {
	return &GetUserOrdersUseCase{orderRepo: orderRepo}
}

func (uc *GetUserOrdersUseCase) Execute(ctx context.Context, userID string) ([]*domain.Order, error) {
	return uc.orderRepo.FindByUserID(ctx, userID)
}
