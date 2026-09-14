package domain

import "context"

type ReviewRepository interface {
	Save(ctx context.Context, review *Review) error

	// FindByProductID - berilgan mahsulot uchun sharhlarni sahifalab qaytaradi (jami son bilan
	// birga, javobda pagination meta-ma'lumot uchun).
	FindByProductID(ctx context.Context, productID string, page, pageSize int) ([]*Review, int64, error)

	FindByUserID(ctx context.Context, userID string) ([]*Review, error)

	// ExistsForUserAndProduct - foydalanuvchi shu mahsulotga allaqachon sharh qoldirganini
	// tekshiradi (bitta foydalanuvchi - bitta mahsulot uchun bitta sharh qoidasi).
	ExistsForUserAndProduct(ctx context.Context, userID, productID string) (bool, error)

	// Delete - sharhni bazadan butunlay o'chiradi (admin moderatsiyasi uchun). Topilmasa
	// ErrReviewNotFound qaytaradi.
	Delete(ctx context.Context, id string) error
}
