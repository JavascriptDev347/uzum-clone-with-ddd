package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
)

// DeleteReviewUseCase - admin uchun: istalgan foydalanuvchi yozgan sharhni butunlay o'chiradi
// (moderatsiya maqsadida). Egalik tekshiruvi yo'q - review kimga tegishli bo'lishidan qat'iy
// nazar, faqat admin roli (RequireRole orqali router darajasida) ushbu use case'ga kirish
// huquqini beradi.
type DeleteReviewUseCase struct {
	reviewRepo domain.ReviewRepository
}

func NewDeleteReviewUseCase(reviewRepo domain.ReviewRepository) *DeleteReviewUseCase {
	return &DeleteReviewUseCase{reviewRepo: reviewRepo}
}

func (uc *DeleteReviewUseCase) Execute(ctx context.Context, reviewID string) error {
	return uc.reviewRepo.Delete(ctx, reviewID)
}
