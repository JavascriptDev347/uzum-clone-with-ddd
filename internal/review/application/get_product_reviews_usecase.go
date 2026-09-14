package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
)

// GetProductReviewsUseCase - bitta mahsulot uchun sharhlar ro'yxatini sahifalab olish (ochiq,
// autentifikatsiya shart emas).
type GetProductReviewsUseCase struct {
	reviewRepo domain.ReviewRepository
}

func NewGetProductReviewsUseCase(reviewRepo domain.ReviewRepository) *GetProductReviewsUseCase {
	return &GetProductReviewsUseCase{reviewRepo: reviewRepo}
}

func (uc *GetProductReviewsUseCase) Execute(ctx context.Context, productID string, page, pageSize int) ([]*domain.Review, int64, error) {
	page, pageSize = NormalizeReviewPagination(page, pageSize)
	return uc.reviewRepo.FindByProductID(ctx, productID, page, pageSize)
}
