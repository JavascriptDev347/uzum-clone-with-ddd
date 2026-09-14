package domain

import "context"

// StockReserver - checkout paytida Catalog'dagi mahsulot zaxirasini kamaytirish/tiklash uchun
// port (Catalog bilan bog'lanish shu interfeys orqali amalga oshiriladi, ordering hech qachon
// catalog'ning repository yoki SQL qatlamiga to'g'ridan-to'g'ri kirmaydi). Implementatsiyasi
// (ordering/infrastructure) Catalog'ning application qatlamiga murojaat qiladi va DB darajasida
// atomik bo'lishi shart - bir vaqtda bir nechta checkout bo'lganda oversell bo'lmasligi uchun.
type StockReserver interface {
	// DecrementStock - berilgan miqdorni mahsulot zaxirasidan ayiradi. Agar yetarli zaxira
	// bo'lmasa ErrInsufficientStock qaytaradi.
	DecrementStock(ctx context.Context, productID string, quantity int) error

	// RestoreStock - muvaffaqiyatsiz checkout'da shu paytgacha ayirilgan miqdorni orqaga
	// qaytaradi (kompensatsiya). Umumiy DB tranzaksiyasi mavjud bo'lmagani uchun ishlatiladi.
	RestoreStock(ctx context.Context, productID string, quantity int) error
}
