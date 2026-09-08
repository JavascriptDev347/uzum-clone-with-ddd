package domain

import (
	"context"
	"errors"
)

var ErrWishlistNotFound = errors.New("wishlist: wishlist topilmadi")

type WishlistRepository interface {
	Save(ctx context.Context, wishlist *Wishlist) error
	FindByUserID(ctx context.Context, userID string) (*Wishlist, error)
}
