package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
)

// GetUserReviewsUseCase - foydalanuvchining o'zi qoldirgan barcha sharhlari ro'yxatini olish.
type GetUserReviewsUseCase struct {
	reviewRepo domain.ReviewRepository
}

func NewGetUserReviewsUseCase(reviewRepo domain.ReviewRepository) *GetUserReviewsUseCase {
	return &GetUserReviewsUseCase{reviewRepo: reviewRepo}
}

func (uc *GetUserReviewsUseCase) Execute(ctx context.Context, userID string) ([]*domain.Review, error) {
	return uc.reviewRepo.FindByUserID(ctx, userID)
}
