package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyProductID           = errors.New("wishlist: product ID bo'sh bo'lishi mumkin emas")
	ErrEmptyWishlistID          = errors.New("wishlist: wishlist ID bo'sh bo'lishi mumkin emas")
	ErrEmptyUserID              = errors.New("wishlist: user ID bo'sh bo'lishi mumkin emas")
	ErrProductAlreadyInWishlist = errors.New("wishlist: product already in wishlist")
	ErrProductNotInWishlist     = errors.New("wishlist: product not in wishlist")
)

type Wishlist struct {
	id        string
	userID    string
	items     []WishlistItem
	createdAt time.Time
}

type WishlistItem struct {
	productID string
	addedAt   time.Time
}

func NewWishlist(id, userID string) (*Wishlist, error) {
	if id == "" {
		return nil, ErrEmptyWishlistID
	}
	if userID == "" {
		return nil, ErrEmptyUserID
	}

	return &Wishlist{
		id:        id,
		userID:    userID,
		items:     []WishlistItem{},
		createdAt: time.Now(),
	}, nil
}

func NewWishlistFromRepository(id, userID string, items []WishlistItem, createdAt time.Time) *Wishlist {
	return &Wishlist{
		id:        id,
		userID:    userID,
		items:     items,
		createdAt: createdAt,
	}
}

func NewWishlistItemFromRepository(productID string, addedAt time.Time) WishlistItem {
	return WishlistItem{
		productID: productID,
		addedAt:   addedAt,
	}
}

func (w *Wishlist) AddItem(productID string) error {
	if productID == "" {
		return ErrEmptyProductID
	}
	if w.HasItem(productID) {
		return ErrProductAlreadyInWishlist
	}
	w.items = append(w.items, WishlistItem{
		productID: productID,
		addedAt:   time.Now(),
	})
	return nil
}

func (w *Wishlist) RemoveItem(productID string) error {
	for i, item := range w.items {
		if item.productID == productID {
			w.items = append(w.items[:i], w.items[i+1:]...)
			return nil
		}
	}
	return ErrProductNotInWishlist
}

func (w *Wishlist) HasItem(productID string) bool {
	for _, item := range w.items {
		if item.productID == productID {
			return true
		}
	}
	return false
}

func (w *Wishlist) ID() string {
	return w.id
}

func (w *Wishlist) UserID() string {
	return w.userID
}

func (w *Wishlist) Items() []WishlistItem {
	return w.items
}

func (w *Wishlist) CreatedAt() time.Time {
	return w.createdAt
}

func (i WishlistItem) ProductID() string {
	return i.productID
}

func (i WishlistItem) AddedAt() time.Time {
	return i.addedAt
}
