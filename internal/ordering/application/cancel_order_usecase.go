package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// CancelOrderUseCase - admin uchun: buyurtmani bekor qiladi. Checkout va CreateManualOrder
// buyurtma yaratilishida zaxirani kamaytirgani uchun, bekor qilishda bu zaxira Catalog'ga
// qaytarilishi kerak (StockReserver.RestoreStock orqali, checkout'dagi rollback bilan bir xil
// port). Faqat Preparing holatidagi buyurtmalar bekor qilinishi mumkinligi mantig'i bu yerda
// takrorlanmaydi - u to'liq Order.Cancel()ga yuklangan (domain aggregate o'z invariantini
// o'zi himoya qiladi).
type CancelOrderUseCase struct {
	orderRepo     domain.OrderRepository
	stockReserver domain.StockReserver
}

func NewCancelOrderUseCase(orderRepo domain.OrderRepository, stockReserver domain.StockReserver) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		orderRepo:     orderRepo,
		stockReserver: stockReserver,
	}
}

func (uc *CancelOrderUseCase) Execute(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	// Zaxirani qaytarishdan oldin tekshiramiz - HandedToCourier/Delivered bo'lsa Cancel()
	// xato qaytaradi va biz hech qanday zaxirani qaytarmasdan chiqamiz.
	if err := order.Cancel(); err != nil {
		return nil, err
	}

	for _, item := range order.Items() {
		if err := uc.stockReserver.RestoreStock(ctx, item.ProductID(), item.Quantity()); err != nil {
			return nil, err
		}
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
