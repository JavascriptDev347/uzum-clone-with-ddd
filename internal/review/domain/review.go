package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MinRating = 1
	MaxRating = 5
)

// ValidateRating - Rating 1 dan 5 gacha oralig'ida ekanligini tekshiradi. NewReview'ning
// o'zi ham shu tekshiruvni qiladi, lekin SubmitReviewUseCase yetkazib berish/duplikat
// tekshiruvlaridan oldin tez orada rad etish uchun buni alohida chaqiradi.
func ValidateRating(rating int) error {
	if rating < MinRating || rating > MaxRating {
		return ErrInvalidRating
	}
	return nil
}

// Review - foydalanuvchi sotib olib, yetkazib berilgan mahsulotga qoldirgan sharhi. OrderID -
// shu sharhni yozishga huquq bergan buyurtma (DeliveredPurchaseChecker orqali aniqlanadi),
// keyinchalik audit uchun saqlanadi.
type Review struct {
	id        uuid.UUID
	userID    string
	productID string
	orderID   string
	rating    int
	comment   string
	createdAt time.Time
}

// NewReview - yangi sharh yaratish uchun. Rating shu yerda tekshiriladi - Review o'z
// invariantini o'zi himoya qiladi.
func NewReview(userID, productID, orderID string, rating int, comment string) (*Review, error) {
	if err := ValidateRating(rating); err != nil {
		return nil, err
	}

	return &Review{
		id:        uuid.New(),
		userID:    userID,
		productID: productID,
		orderID:   orderID,
		rating:    rating,
		comment:   comment,
		createdAt: time.Now(),
	}, nil
}

// NewReviewFromRepository - saqlangan sharhni bazadan qayta tiklash uchun, domen
// konstruktoridan ajratilgan (invariantlar bazadan o'qishda qayta tekshirilmaydi).
func NewReviewFromRepository(
	id uuid.UUID,
	userID, productID, orderID string,
	rating int,
	comment string,
	createdAt time.Time,
) *Review {
	return &Review{
		id:        id,
		userID:    userID,
		productID: productID,
		orderID:   orderID,
		rating:    rating,
		comment:   comment,
		createdAt: createdAt,
	}
}

func (r *Review) ID() uuid.UUID        { return r.id }
func (r *Review) UserID() string       { return r.userID }
func (r *Review) ProductID() string    { return r.productID }
func (r *Review) OrderID() string      { return r.orderID }
func (r *Review) Rating() int          { return r.rating }
func (r *Review) Comment() string      { return r.comment }
func (r *Review) CreatedAt() time.Time { return r.createdAt }
