package domain

import (
	"errors"
	"testing"
)

func TestNewReview_RatingBounds(t *testing.T) {
	tests := []struct {
		name    string
		rating  int
		wantErr error
	}{
		{name: "min valid rating", rating: 1, wantErr: nil},
		{name: "max valid rating", rating: 5, wantErr: nil},
		{name: "mid valid rating", rating: 3, wantErr: nil},
		{name: "zero rejected", rating: 0, wantErr: ErrInvalidRating},
		{name: "negative rejected", rating: -1, wantErr: ErrInvalidRating},
		{name: "above max rejected", rating: 6, wantErr: ErrInvalidRating},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			review, err := NewReview("user-1", "product-1", "order-1", tt.rating, "")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}
			if review.Rating() != tt.rating {
				t.Fatalf("expected rating %d, got %d", tt.rating, review.Rating())
			}
		})
	}
}

func TestNewReview_FieldsAreSet(t *testing.T) {
	review, err := NewReview("user-1", "product-1", "order-1", 4, "Yaxshi mahsulot")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if review.UserID() != "user-1" {
		t.Fatalf("expected UserID user-1, got %v", review.UserID())
	}
	if review.ProductID() != "product-1" {
		t.Fatalf("expected ProductID product-1, got %v", review.ProductID())
	}
	if review.OrderID() != "order-1" {
		t.Fatalf("expected OrderID order-1, got %v", review.OrderID())
	}
	if review.Comment() != "Yaxshi mahsulot" {
		t.Fatalf("expected comment to be preserved, got %v", review.Comment())
	}
	if review.ID().String() == "" {
		t.Fatal("expected a generated ID")
	}
}

func TestValidateRating(t *testing.T) {
	tests := []struct {
		rating  int
		wantErr error
	}{
		{rating: 1, wantErr: nil},
		{rating: 5, wantErr: nil},
		{rating: 0, wantErr: ErrInvalidRating},
		{rating: 6, wantErr: ErrInvalidRating},
	}

	for _, tt := range tests {
		if err := ValidateRating(tt.rating); !errors.Is(err, tt.wantErr) {
			t.Fatalf("ValidateRating(%d): expected %v, got %v", tt.rating, tt.wantErr, err)
		}
	}
}
