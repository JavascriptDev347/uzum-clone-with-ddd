package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// GetOrderByIDUseCase - bitta buyurtmani olish, egalik tekshiruvi bilan: admin bo'lmasa,
// faqat o'ziga tegishli buyurtmani ko'ra oladi. isAdmin hozircha parametr sifatida qabul
// qilinadi - haqiqiy rol tekshiruvi keyingi qadamda HTTP qatlamida ulanadi.
type GetOrderByIDUseCase struct {
	orderRepo domain.OrderRepository
}

func NewGetOrderByIDUseCase(orderRepo domain.OrderRepository) *GetOrderByIDUseCase {
	return &GetOrderByIDUseCase{orderRepo: orderRepo}
}

func (uc *GetOrderByIDUseCase) Execute(ctx context.Context, orderID, requestingUserID string, isAdmin bool) (*domain.Order, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	if !isAdmin {
		// UserID nil bo'lsa (admin tomonidan qo'lda yaratilgan offline buyurtma) - hech qanday
		// mijoz hisobi unga tegishli emas, shuning uchun admin bo'lmagan so'rovchi ko'ra olmaydi.
		if order.UserID() == nil || *order.UserID() != requestingUserID {
			return nil, domain.ErrOrderAccessDenied
		}
	}

	return order, nil
}
