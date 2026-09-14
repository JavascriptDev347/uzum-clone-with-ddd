package http

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
)

// SubmitReviewRequest - sharh qoldirish so'rovining tanasi. Comment ixtiyoriy.
type SubmitReviewRequest struct {
	ProductID string `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment,omitempty"`
}

// ReviewResponse - sharhning HTTP javob shakli. domain.Review to'g'ridan-to'g'ri
// serializatsiya qilinmaydi (private field'lari bor), shu sabab bu alohida DTO orqali map
// qilinadi.
type ReviewResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	OrderID   string    `json:"order_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func ToReviewResponse(review *domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:        review.ID().String(),
		UserID:    review.UserID(),
		ProductID: review.ProductID(),
		OrderID:   review.OrderID(),
		Rating:    review.Rating(),
		Comment:   review.Comment(),
		CreatedAt: review.CreatedAt(),
	}
}

func ToReviewResponses(reviews []*domain.Review) []ReviewResponse {
	responses := make([]ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		responses = append(responses, ToReviewResponse(r))
	}
	return responses
}
