package application

import (
	"context"
	"errors"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type CartItemView struct {
	ProductID     string   `json:"product_id"`
	ProductName   string   `json:"product_name"`
	UnitPrice     float64  `json:"unit_price"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
	Currency      string   `json:"currency"`
	Quantity      int      `json:"quantity"`
	Subtotal      float64  `json:"subtotal"` // discount bo'lsa discount narxidan, aks holda asl narxdan hisoblanadi
	Available     bool     `json:"available"`
}

type CartView struct {
	Items      []CartItemView `json:"items"`
	TotalItems int            `json:"total_items"`
	TotalPrice float64        `json:"total_price"` // barcha item'lar subtotal'larining yig'indisi
}

type GetCartUseCase struct {
	cartRepo    domain.CartRepository
	productRepo catalog.ProductRepository
}

func NewGetCartUseCase(cartRepo domain.CartRepository, productRepo catalog.ProductRepository) *GetCartUseCase {
	return &GetCartUseCase{cartRepo: cartRepo, productRepo: productRepo}
}

func (uc *GetCartUseCase) Execute(ctx context.Context, userID string) (*CartView, error) {
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cart mavjud emas — bo'sh CartView qaytaramiz, xato emas
	if cart == nil {
		return &CartView{Items: []CartItemView{}, TotalItems: 0}, nil
	}

	views := make([]CartItemView, 0, len(cart.Items()))
	var totalPrice float64
	for _, item := range cart.Items() {
		product, err := uc.productRepo.FindByID(ctx, item.ProductID())
		if err != nil && !errors.Is(err, catalog.ErrProductNotFound) {
			// TODO: bu haqiqiy infra xato (DB ulanish va h.k.) — butun so'rovni to'xtatamiz
			return nil, err
		}

		if err != nil || product == nil {
			// Mahsulot o'chirilgan yoki topilmagan — item ko'rsatiladi, lekin "mavjud emas" deb belgilanadi
			views = append(views, CartItemView{
				ProductID: item.ProductID(),
				Quantity:  item.Quantity(),
				Available: false,
			})
			continue
		}

		unitPrice := product.Price()
		effectiveAmount := unitPrice.Amount()

		var discountPrice *float64
		if discount := product.DiscountPrice(); discount != nil {
			amount := discount.Amount()
			discountPrice = &amount
			effectiveAmount = amount // discount bo'lsa, hisob-kitob discount narxidan olib boriladi
		}

		subtotal := effectiveAmount * float64(item.Quantity())
		totalPrice += subtotal

		views = append(views, CartItemView{
			ProductID:     item.ProductID(),
			ProductName:   product.NameUz(),
			UnitPrice:     unitPrice.Amount(),
			DiscountPrice: discountPrice,
			Currency:      unitPrice.Currency(),
			Quantity:      item.Quantity(),
			Subtotal:      subtotal,
			Available:     true,
		})
	}

	return &CartView{
		Items:      views,
		TotalItems: cart.TotalItems(),
		TotalPrice: totalPrice,
	}, nil
}
