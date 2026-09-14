package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
)

// fakeReviewRepository - domain.ReviewRepository ning xotiradagi implementatsiyasi, testlar
// uchun.
type fakeReviewRepository struct {
	reviews   map[string]*domain.Review // key: userID + "|" + productID
	saveErr   error
	existsErr error
	saveCalls int
}

func newFakeReviewRepository() *fakeReviewRepository {
	return &fakeReviewRepository{reviews: make(map[string]*domain.Review)}
}

func reviewKey(userID, productID string) string {
	return userID + "|" + productID
}

func (f *fakeReviewRepository) Save(ctx context.Context, review *domain.Review) error {
	f.saveCalls++
	if f.saveErr != nil {
		return f.saveErr
	}
	f.reviews[reviewKey(review.UserID(), review.ProductID())] = review
	return nil
}

func (f *fakeReviewRepository) FindByProductID(ctx context.Context, productID string, page, pageSize int) ([]*domain.Review, int64, error) {
	result := make([]*domain.Review, 0)
	for _, r := range f.reviews {
		if r.ProductID() == productID {
			result = append(result, r)
		}
	}
	return result, int64(len(result)), nil
}

func (f *fakeReviewRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Review, error) {
	result := make([]*domain.Review, 0)
	for _, r := range f.reviews {
		if r.UserID() == userID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (f *fakeReviewRepository) ExistsForUserAndProduct(ctx context.Context, userID, productID string) (bool, error) {
	if f.existsErr != nil {
		return false, f.existsErr
	}
	_, ok := f.reviews[reviewKey(userID, productID)]
	return ok, nil
}

// fakePurchaseChecker - domain.DeliveredPurchaseChecker ning test uchun soxta
// implementatsiyasi.
type fakePurchaseChecker struct {
	deliveredOrderID string
	eligible         bool
	err              error
	calls            int
}

func (f *fakePurchaseChecker) HasDeliveredPurchase(ctx context.Context, userID, productID string) (string, bool, error) {
	f.calls++
	if f.err != nil {
		return "", false, f.err
	}
	if !f.eligible {
		return "", false, nil
	}
	return f.deliveredOrderID, true, nil
}

func TestSubmitReviewUseCase_Execute_EligibleFirstReviewSucceeds(t *testing.T) {
	reviewRepo := newFakeReviewRepository()
	purchaseChecker := &fakePurchaseChecker{eligible: true, deliveredOrderID: "order-1"}

	uc := application.NewSubmitReviewUseCase(reviewRepo, purchaseChecker)

	review, err := uc.Execute(context.Background(), application.SubmitReviewInput{
		UserID:    "user-1",
		ProductID: "product-1",
		Rating:    5,
		Comment:   "Ajoyib!",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if review.OrderID() != "order-1" {
		t.Fatalf("expected OrderID from purchase checker, got %v", review.OrderID())
	}
	if review.Rating() != 5 {
		t.Fatalf("expected rating 5, got %d", review.Rating())
	}
	if reviewRepo.saveCalls != 1 {
		t.Fatalf("expected exactly 1 save call, got %d", reviewRepo.saveCalls)
	}
}

func TestSubmitReviewUseCase_Execute_IneligibleRejectedBeforeAnyDBWrite(t *testing.T) {
	reviewRepo := newFakeReviewRepository()
	purchaseChecker := &fakePurchaseChecker{eligible: false}

	uc := application.NewSubmitReviewUseCase(reviewRepo, purchaseChecker)

	_, err := uc.Execute(context.Background(), application.SubmitReviewInput{
		UserID:    "user-1",
		ProductID: "product-1",
		Rating:    4,
	})
	if !errors.Is(err, domain.ErrNotEligibleForReview) {
		t.Fatalf("expected ErrNotEligibleForReview, got %v", err)
	}
	if reviewRepo.saveCalls != 0 {
		t.Fatalf("expected no save calls for an ineligible user, got %d", reviewRepo.saveCalls)
	}
}

func TestSubmitReviewUseCase_Execute_DuplicateReviewRejected(t *testing.T) {
	reviewRepo := newFakeReviewRepository()
	purchaseChecker := &fakePurchaseChecker{eligible: true, deliveredOrderID: "order-1"}

	uc := application.NewSubmitReviewUseCase(reviewRepo, purchaseChecker)

	// Birinchi sharh muvaffaqiyatli.
	if _, err := uc.Execute(context.Background(), application.SubmitReviewInput{
		UserID:    "user-1",
		ProductID: "product-1",
		Rating:    5,
	}); err != nil {
		t.Fatalf("unexpected error on first review: %v", err)
	}

	// Xuddi shu foydalanuvchi, xuddi shu mahsulot uchun ikkinchi urinish rad etiladi.
	_, err := uc.Execute(context.Background(), application.SubmitReviewInput{
		UserID:    "user-1",
		ProductID: "product-1",
		Rating:    3,
	})
	if !errors.Is(err, domain.ErrReviewAlreadyExists) {
		t.Fatalf("expected ErrReviewAlreadyExists, got %v", err)
	}
	if reviewRepo.saveCalls != 1 {
		t.Fatalf("expected save to have been called only once (for the first review), got %d", reviewRepo.saveCalls)
	}
}

func TestSubmitReviewUseCase_Execute_InvalidRatingRejectedBeforeEligibilityCheck(t *testing.T) {
	reviewRepo := newFakeReviewRepository()
	purchaseChecker := &fakePurchaseChecker{eligible: true, deliveredOrderID: "order-1"}

	uc := application.NewSubmitReviewUseCase(reviewRepo, purchaseChecker)

	_, err := uc.Execute(context.Background(), application.SubmitReviewInput{
		UserID:    "user-1",
		ProductID: "product-1",
		Rating:    0,
	})
	if !errors.Is(err, domain.ErrInvalidRating) {
		t.Fatalf("expected ErrInvalidRating, got %v", err)
	}
	if purchaseChecker.calls != 0 {
		t.Fatalf("expected no eligibility check for an invalid rating, got %d calls", purchaseChecker.calls)
	}
	if reviewRepo.saveCalls != 0 {
		t.Fatalf("expected no save calls for an invalid rating, got %d", reviewRepo.saveCalls)
	}
}
