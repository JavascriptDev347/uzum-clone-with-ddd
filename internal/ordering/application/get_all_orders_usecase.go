package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// GetAllOrdersUseCase - admin buyurtmalar navbati uchun: barcha buyurtmalar ro'yxati.
type GetAllOrdersUseCase struct {
	orderRepo domain.OrderRepository
}

func NewGetAllOrdersUseCase(orderRepo domain.OrderRepository) *GetAllOrdersUseCase {
	return &GetAllOrdersUseCase{orderRepo: orderRepo}
}

func (uc *GetAllOrdersUseCase) Execute(ctx context.Context) ([]*domain.Order, error) {
	return uc.orderRepo.FindAll(ctx)
}
