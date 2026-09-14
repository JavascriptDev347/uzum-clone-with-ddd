package domain

import "errors"

var (
	ErrInvalidRating        = errors.New("review: reyting 1 dan 5 gacha bo'lishi kerak")
	ErrNotEligibleForReview = errors.New("review: bu mahsulot uchun sharh qoldirish huquqingiz yo'q")
	ErrReviewAlreadyExists  = errors.New("review: bu mahsulot uchun sharh allaqachon qoldirilgan")
)
