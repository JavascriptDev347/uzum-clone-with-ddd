package domain

import "errors"

var (
	ErrEmptyAddress              = errors.New("ordering: manzil bo'sh bo'lishi mumkin emas")
	ErrInvalidPhone              = errors.New("ordering: telefon raqami formati noto'g'ri")
	ErrEmptyOrderItems           = errors.New("ordering: buyurtmada kamida bitta mahsulot bo'lishi kerak")
	ErrEmptyOrderItemProductID   = errors.New("ordering: order item uchun product ID bo'sh bo'lishi mumkin emas")
	ErrEmptyOrderItemName        = errors.New("ordering: order item uchun mahsulot nomi bo'sh bo'lishi mumkin emas")
	ErrInvalidQuantity           = errors.New("ordering: miqdor noldan katta bo'lishi kerak")
	ErrInvalidDeliveryStatus     = errors.New("ordering: yetkazib berish holati noto'g'ri")
	ErrInvalidDeliveryTransition = errors.New("ordering: yetkazib berish holatini orqaga qaytarib bo'lmaydi")
	ErrOrderNotFound             = errors.New("ordering: order topilmadi")
	ErrOrderAccessDenied         = errors.New("ordering: bu buyurtmani ko'rish huquqingiz yo'q")
	ErrInsufficientStock         = errors.New("ordering: mahsulot uchun yetarli miqdorda zaxira yo'q")
	ErrInvalidPaymentStatus      = errors.New("ordering: to'lov holati noto'g'ri")
	ErrOrderAlreadyShipped       = errors.New("ordering: courier'ga topshirilgan yoki yetkazilgan buyurtmani bekor qilib bo'lmaydi")
)
