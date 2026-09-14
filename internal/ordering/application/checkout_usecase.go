package application

import (
	"context"

	cartapp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/application"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// CartReader - cart context'idan checkout uchun kerakli joriy cart holatini o'qish uchun
// chegara (ACL). cart/application'dagi mavjud GetCartUseCase buni qondiradi (Execute imzosi
// mos keladi) - ordering hech qachon cart'ning domain yoki repository qatlamiga to'g'ridan-to'g'ri
// kirmaydi.
type CartReader interface {
	Execute(ctx context.Context, userID string) (*cartapp.CartView, error)
}

// CartClearer - checkout muvaffaqiyatli yakunlangandan keyin cart'ni bo'shatish uchun chegara.
// cart/application'dagi ClearCartUseCase buni qondiradi.
type CartClearer interface {
	Execute(ctx context.Context, userID string) error
}

type CheckoutUseCase struct {
	orderRepo     domain.OrderRepository
	cartReader    CartReader
	cartClearer   CartClearer
	productRepo   catalog.ProductRepository
	stockReserver domain.StockReserver
}

func NewCheckoutUseCase(
	orderRepo domain.OrderRepository,
	cartReader CartReader,
	cartClearer CartClearer,
	productRepo catalog.ProductRepository,
	stockReserver domain.StockReserver,
) *CheckoutUseCase {
	return &CheckoutUseCase{
		orderRepo:     orderRepo,
		cartReader:    cartReader,
		cartClearer:   cartClearer,
		productRepo:   productRepo,
		stockReserver: stockReserver,
	}
}

type CheckoutInput struct {
	UserID  string
	Address string
	Phone   string
	Note    string
}

func (uc *CheckoutUseCase) Execute(ctx context.Context, in CheckoutInput) (*domain.Order, error) {
	customerInfo, err := domain.NewCustomerInfo(in.Address, in.Phone, in.Note)
	if err != nil {
		return nil, err
	}

	cart, err := uc.cartReader.Execute(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, domain.ErrEmptyOrderItems
	}

	requests := make([]productQuantity, 0, len(cart.Items))
	for _, cartItem := range cart.Items {
		if !cartItem.Available {
			return nil, catalog.ErrProductNotFound
		}
		requests = append(requests, productQuantity{ProductID: cartItem.ProductID, Quantity: cartItem.Quantity})
	}

	// Har bir cart qatori uchun Catalog'dan eng yangi mahsulot ma'lumotini olamiz - nom va
	// narx shu yerda "surat" (snapshot) sifatida qotib qoladi (buildOrderItemSnapshots orqali -
	// CreateManualOrder bilan baham ko'riladi, mantiq takrorlanmaydi).
	items, err := buildOrderItemSnapshots(ctx, uc.productRepo, requests)
	if err != nil {
		return nil, err
	}

	// Zaxirani buyurtma yaratilishidan OLDIN, atomik ravishda kamaytiramiz (reserveOrderStock
	// orqali - CreateManualOrder bilan baham ko'riladi) - bir vaqtda bir nechta checkout
	// bo'lganda oversell bo'lmasligi uchun.
	if err := reserveOrderStock(ctx, uc.stockReserver, items); err != nil {
		return nil, err
	}

	order, err := domain.NewOrder(&in.UserID, customerInfo, items)
	if err != nil {
		releaseOrderStock(ctx, uc.stockReserver, items)
		return nil, err
	}

	if err := uc.orderRepo.Save(ctx, order); err != nil {
		releaseOrderStock(ctx, uc.stockReserver, items)
		return nil, err
	}

	// Buyurtma allaqachon saqlandi, shuning uchun bu yerdagi xatolik buyurtmani bekor
	// qilmaydi - lekin chaqiruvchiga xabar qilinishi uchun qaytariladi.
	if err := uc.cartClearer.Execute(ctx, in.UserID); err != nil {
		return order, err
	}

	return order, nil
}
