package domain

import "context"

// DeliveredPurchaseChecker - ordering context'iga ACL chegarasi: review'ning domain va
// application qatlamlari ordering/domain yoki ordering/infrastructure'ni hech qachon
// to'g'ridan-to'g'ri import qilmaydi - faqat shu port orqali murojaat qiladi.
// Implementatsiyasi (review/infrastructure) ordering'ning application qatlamiga murojaat
// qiladi, hech qachon uning bazasiga to'g'ridan-to'g'ri emas.
type DeliveredPurchaseChecker interface {
	// HasDeliveredPurchase - foydalanuvchi berilgan mahsulotni o'z ichiga olgan va allaqachon
	// yetkazib berilgan (Delivered) buyurtmaga ega bo'lsa, shu buyurtmaning ID'sini qaytaradi.
	HasDeliveredPurchase(ctx context.Context, userID, productID string) (orderID string, ok bool, err error)
}
