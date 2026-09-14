package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// HasDeliveredProductUseCase - boshqa bounded context'lar (masalan review) uchun fokuslangan,
// faqat-o'qish uchun so'rov: foydalanuvchi berilgan mahsulotni o'z ichiga olgan va allaqachon
// yetkazib berilgan (Delivered) buyurtmaga egami. GetUserOrdersUseCase'dan foydalanish
// chaqiruvchiga kerak bo'lmagan barcha buyurtma tafsilotlarini (customerInfo, boshqa itemlar
// va h.k.) oshkor qilar edi va "Delivered + mahsulot mavjud" tekshiruvini chaqiruvchi
// tomonga yuklardi - bu ordering'ning o'z domenini bilishi kerak bo'lgan mantiq, shuning
// uchun bu yerda, ordering'ning application qatlamida qoladi.
type HasDeliveredProductUseCase struct {
	orderRepo domain.OrderRepository
}

func NewHasDeliveredProductUseCase(orderRepo domain.OrderRepository) *HasDeliveredProductUseCase {
	return &HasDeliveredProductUseCase{orderRepo: orderRepo}
}

// Execute - foydalanuvchining Delivered holatidagi buyurtmalari orasidan berilgan mahsulotni
// o'z ichiga olgan birinchisini topadi va shu buyurtma ID'sini qaytaradi.
func (uc *HasDeliveredProductUseCase) Execute(ctx context.Context, userID, productID string) (orderID string, ok bool, err error) {
	orders, err := uc.orderRepo.FindByUserID(ctx, userID)
	if err != nil {
		return "", false, err
	}

	for _, order := range orders {
		if order.DeliveryStatus() != domain.DeliveryStatusDelivered {
			continue
		}
		for _, item := range order.Items() {
			if item.ProductID() == productID {
				return order.ID().String(), true, nil
			}
		}
	}

	return "", false, nil
}
