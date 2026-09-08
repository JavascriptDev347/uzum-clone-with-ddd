package application

import (
	"context"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/domain"
)

type WishlistItemOutput struct {
	ProductID string `json:"product_id"`
	AddedAt   string `json:"added_at"` // RFC3339 formatida, time.Time emas
}

type GetWishlistOutput struct {
	UserID string               `json:"user_id"`
	Items  []WishlistItemOutput `json:"items"`
}

type GetWishlistUseCase struct {
	repo domain.WishlistRepository
}

func NewGetWishlistUseCase(repo domain.WishlistRepository) *GetWishlistUseCase {
	return &GetWishlistUseCase{repo: repo}
}

func (uc *GetWishlistUseCase) Execute(ctx context.Context, userID string) (*GetWishlistOutput, error) {
	wishlist, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrWishlistNotFound) {
			return &GetWishlistOutput{
				UserID: userID,
				Items:  make([]WishlistItemOutput, 0),
			}, nil
		}
		return nil, err
	}

	output := &GetWishlistOutput{
		UserID: userID,
		Items:  make([]WishlistItemOutput, 0),
	}
	for _, item := range wishlist.Items() {
		output.Items = append(output.Items, WishlistItemOutput{
			ProductID: item.ProductID(),
			AddedAt:   item.AddedAt().Format(time.RFC3339),
		})
	}
	return output, nil
}
