package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/domain"
)

type RemoveFromWishlistInput struct {
	UserID    string
	ProductID string
}

type RemoveFromWishlistUseCase struct {
	repo domain.WishlistRepository
}

func NewRemoveFromWishlistUseCase(repo domain.WishlistRepository) *RemoveFromWishlistUseCase {
	return &RemoveFromWishlistUseCase{repo: repo}
}

func (uc *RemoveFromWishlistUseCase) Execute(ctx context.Context, input RemoveFromWishlistInput) error {
	wishlist, err := uc.repo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if err := wishlist.RemoveItem(input.ProductID); err != nil {
		return err
	}

	return uc.repo.Save(ctx, wishlist)
}
