package application

import (
	"context"
	"errors"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/domain"
	"github.com/google/uuid"
)

type AddToWishlistInput struct {
	UserID    string
	ProductID string
}

type AddToWishlistUseCase struct {
	repo domain.WishlistRepository
}

func NewAddToWishlistUseCase(repo domain.WishlistRepository) *AddToWishlistUseCase {
	return &AddToWishlistUseCase{repo: repo}
}

func (uc *AddToWishlistUseCase) Execute(ctx context.Context, input AddToWishlistInput) error {

	wishlist, err := uc.repo.FindByUserID(ctx, input.UserID)
	if err != nil {
		if !errors.Is(err, domain.ErrWishlistNotFound) {
			return err
		}
		wishlist, err = domain.NewWishlist(uuid.New().String(), input.UserID)
		if err != nil {
			return err
		}
	}

	if err := wishlist.AddItem(input.ProductID); err != nil {
		return err
	}

	return uc.repo.Save(ctx, wishlist)
}
