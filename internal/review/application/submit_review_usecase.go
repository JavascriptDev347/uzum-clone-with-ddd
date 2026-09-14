package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
)

// SubmitReviewInput - sharh qoldirish uchun kirish ma'lumotlari.
type SubmitReviewInput struct {
	UserID    string
	ProductID string
	Rating    int
	Comment   string
}

// SubmitReviewUseCase - foydalanuvchi mahsulotga sharh qoldiradi. Faqat shu mahsulotni sotib
// olib, yetkazib berilgan (Delivered) buyurtmasi bo'lgan foydalanuvchilar sharh qoldira oladi
// (DeliveredPurchaseChecker orqali tekshiriladi - ordering context'iga ACL chegarasi), va har
// bir foydalanuvchi bitta mahsulotga faqat bitta marta sharh qoldirishi mumkin.
type SubmitReviewUseCase struct {
	reviewRepo      domain.ReviewRepository
	purchaseChecker domain.DeliveredPurchaseChecker
}

func NewSubmitReviewUseCase(reviewRepo domain.ReviewRepository, purchaseChecker domain.DeliveredPurchaseChecker) *SubmitReviewUseCase {
	return &SubmitReviewUseCase{
		reviewRepo:      reviewRepo,
		purchaseChecker: purchaseChecker,
	}
}

func (uc *SubmitReviewUseCase) Execute(ctx context.Context, in SubmitReviewInput) (*domain.Review, error) {
	// Eng arzon tekshiruv birinchi: ordering'ga (tarmoq/DB) yoki reviewRepo'ga murojaat
	// qilishdan oldin noto'g'ri Rating'ni rad etamiz.
	if err := domain.ValidateRating(in.Rating); err != nil {
		return nil, err
	}

	orderID, ok, err := uc.purchaseChecker.HasDeliveredPurchase(ctx, in.UserID, in.ProductID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrNotEligibleForReview
	}

	exists, err := uc.reviewRepo.ExistsForUserAndProduct(ctx, in.UserID, in.ProductID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrReviewAlreadyExists
	}

	review, err := domain.NewReview(in.UserID, in.ProductID, orderID, in.Rating, in.Comment)
	if err != nil {
		return nil, err
	}

	if err := uc.reviewRepo.Save(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}
