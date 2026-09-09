package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

// fakeCartRepository - domain.CartRepository ning xotiradagi implementatsiyasi, testlar uchun.
type fakeCartRepository struct {
	carts map[string]*domain.Cart // userID -> cart
}

func newFakeCartRepository() *fakeCartRepository {
	return &fakeCartRepository{carts: make(map[string]*domain.Cart)}
}

func (f *fakeCartRepository) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	c, ok := f.carts[userID]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (f *fakeCartRepository) Save(ctx context.Context, c *domain.Cart) error {
	f.carts[c.UserID()] = c
	return nil
}

// fakeProductRepository - catalog.ProductRepository ning xotiradagi implementatsiyasi, testlar uchun.
type fakeProductRepository struct {
	products map[string]*catalog.Product
}

func newFakeProductRepository() *fakeProductRepository {
	return &fakeProductRepository{products: make(map[string]*catalog.Product)}
}

func (f *fakeProductRepository) Save(ctx context.Context, product *catalog.Product) error {
	f.products[product.ID()] = product
	return nil
}

func (f *fakeProductRepository) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	p, ok := f.products[id]
	if !ok {
		return nil, catalog.ErrProductNotFound
	}
	return p, nil
}

func (f *fakeProductRepository) FindBySlug(ctx context.Context, slug string) (*catalog.Product, error) {
	return nil, catalog.ErrProductNotFound
}

func (f *fakeProductRepository) FindAll(ctx context.Context, search, categoryID string, page, pageSize int) ([]*catalog.Product, int64, error) {
	return nil, 0, nil
}

func (f *fakeProductRepository) FindAllIncludingDeleted(ctx context.Context, search, categoryID string, page, pageSize int) ([]*catalog.Product, int64, error) {
	return nil, 0, nil
}

func (f *fakeProductRepository) Update(ctx context.Context, product *catalog.Product) error {
	f.products[product.ID()] = product
	return nil
}

func (f *fakeProductRepository) SoftDelete(ctx context.Context, id string) error {
	delete(f.products, id)
	return nil
}

func mustNewProduct(t *testing.T, id, nameUz string, amount float64, currency string, stock int) *catalog.Product {
	t.Helper()
	price, err := catalog.NewMoney(amount, currency)
	if err != nil {
		t.Fatalf("unexpected error creating money: %v", err)
	}
	product, err := catalog.NewProduct(catalog.NewProductParams{
		ID:         id,
		NameUz:     nameUz,
		NameEng:    nameUz,
		NameRu:     nameUz,
		CategoryID: "category-1",
		Price:      price,
		Slug:       nameUz,
		Stock:      stock,
	})
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}
	return product
}

func mustNewProductWithDiscount(t *testing.T, id, nameUz string, amount, discountAmount float64, currency string, stock int) *catalog.Product {
	t.Helper()
	price, err := catalog.NewMoney(amount, currency)
	if err != nil {
		t.Fatalf("unexpected error creating money: %v", err)
	}
	discountPrice, err := catalog.NewMoney(discountAmount, currency)
	if err != nil {
		t.Fatalf("unexpected error creating discount money: %v", err)
	}
	product, err := catalog.NewProduct(catalog.NewProductParams{
		ID:            id,
		NameUz:        nameUz,
		NameEng:       nameUz,
		NameRu:        nameUz,
		CategoryID:    "category-1",
		Price:         price,
		DiscountPrice: &discountPrice,
		Slug:          nameUz,
		Stock:         stock,
	})
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}
	return product
}

func TestAddItemUseCase_Execute(t *testing.T) {
	cartRepo := newFakeCartRepository()
	productRepo := newFakeProductRepository()
	productRepo.products["product-1"] = mustNewProduct(t, "product-1", "Test product", 1000, "UZS", 5)

	uc := application.NewAddItemUseCase(cartRepo, productRepo)

	t.Run("lazy-creates cart and adds item", func(t *testing.T) {
		err := uc.Execute(context.Background(), application.AddItemInput{
			UserID:    "user-1",
			ProductID: "product-1",
			Quantity:  2,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cart, _ := cartRepo.GetByUserID(context.Background(), "user-1")
		if cart == nil {
			t.Fatal("expected cart to be created")
		}
		if cart.TotalItems() != 2 {
			t.Fatalf("expected total 2, got %d", cart.TotalItems())
		}
	})

	t.Run("unknown product returns ErrProductNotFound", func(t *testing.T) {
		err := uc.Execute(context.Background(), application.AddItemInput{
			UserID:    "user-1",
			ProductID: "missing-product",
			Quantity:  1,
		})
		if !errors.Is(err, catalog.ErrProductNotFound) {
			t.Fatalf("expected ErrProductNotFound, got %v", err)
		}
	})

	t.Run("exceeding stock returns ErrInsufficientStock", func(t *testing.T) {
		// cart already has 2 units of product-1 (stock: 5) — adding 4 more makes 6, over stock
		err := uc.Execute(context.Background(), application.AddItemInput{
			UserID:    "user-1",
			ProductID: "product-1",
			Quantity:  4,
		})
		if !errors.Is(err, domain.ErrInsufficientStock) {
			t.Fatalf("expected ErrInsufficientStock, got %v", err)
		}

		// cart quantity must stay unchanged
		cart, _ := cartRepo.GetByUserID(context.Background(), "user-1")
		if cart.TotalItems() != 2 {
			t.Fatalf("expected total to remain 2, got %d", cart.TotalItems())
		}
	})
}

func TestRemoveItemUseCase_Execute(t *testing.T) {
	cartRepo := newFakeCartRepository()
	c := domain.NewCart("user-1")
	_ = c.AddItem("product-1", 2)
	_ = cartRepo.Save(context.Background(), c)

	uc := application.NewRemoveItemUseCase(cartRepo)

	t.Run("removes existing item", func(t *testing.T) {
		if err := uc.Execute(context.Background(), "user-1", "product-1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cart, _ := cartRepo.GetByUserID(context.Background(), "user-1")
		if cart.TotalItems() != 0 {
			t.Fatalf("expected empty cart, got total %d", cart.TotalItems())
		}
	})

	t.Run("cart not found returns ErrItemNotFound", func(t *testing.T) {
		err := uc.Execute(context.Background(), "user-without-cart", "product-1")
		if !errors.Is(err, domain.ErrItemNotFound) {
			t.Fatalf("expected ErrItemNotFound, got %v", err)
		}
	})
}

func TestUpdateItemQuantityUseCase_Execute(t *testing.T) {
	cartRepo := newFakeCartRepository()
	productRepo := newFakeProductRepository()
	productRepo.products["product-1"] = mustNewProduct(t, "product-1", "Test product", 1000, "UZS", 5)

	c := domain.NewCart("user-1")
	_ = c.AddItem("product-1", 2)
	_ = cartRepo.Save(context.Background(), c)

	uc := application.NewUpdateItemQuantityUseCase(cartRepo, productRepo)

	t.Run("updates quantity", func(t *testing.T) {
		err := uc.Execute(context.Background(), application.UpdateItemQuantityInput{
			UserID:    "user-1",
			ProductID: "product-1",
			Quantity:  5,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cart, _ := cartRepo.GetByUserID(context.Background(), "user-1")
		if cart.TotalItems() != 5 {
			t.Fatalf("expected total 5, got %d", cart.TotalItems())
		}
	})

	t.Run("exceeding stock returns ErrInsufficientStock", func(t *testing.T) {
		err := uc.Execute(context.Background(), application.UpdateItemQuantityInput{
			UserID:    "user-1",
			ProductID: "product-1",
			Quantity:  6,
		})
		if !errors.Is(err, domain.ErrInsufficientStock) {
			t.Fatalf("expected ErrInsufficientStock, got %v", err)
		}
		cart, _ := cartRepo.GetByUserID(context.Background(), "user-1")
		if cart.TotalItems() != 5 {
			t.Fatalf("expected total to remain 5, got %d", cart.TotalItems())
		}
	})

	t.Run("zero quantity removes item", func(t *testing.T) {
		err := uc.Execute(context.Background(), application.UpdateItemQuantityInput{
			UserID:    "user-1",
			ProductID: "product-1",
			Quantity:  0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cart, _ := cartRepo.GetByUserID(context.Background(), "user-1")
		if cart.TotalItems() != 0 {
			t.Fatalf("expected empty cart, got total %d", cart.TotalItems())
		}
	})
}

func TestGetCartUseCase_Execute(t *testing.T) {
	cartRepo := newFakeCartRepository()
	productRepo := newFakeProductRepository()
	productRepo.products["product-1"] = mustNewProduct(t, "product-1", "Test product", 1000, "UZS", 10)

	uc := application.NewGetCartUseCase(cartRepo, productRepo)

	t.Run("no cart returns empty view", func(t *testing.T) {
		view, err := uc.Execute(context.Background(), "user-without-cart")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(view.Items) != 0 || view.TotalItems != 0 {
			t.Fatalf("expected empty view, got %+v", view)
		}
	})

	t.Run("returns items with product details, no discount", func(t *testing.T) {
		c := domain.NewCart("user-1")
		_ = c.AddItem("product-1", 3)
		_ = cartRepo.Save(context.Background(), c)

		view, err := uc.Execute(context.Background(), "user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if view.TotalItems != 3 {
			t.Fatalf("expected total 3, got %d", view.TotalItems)
		}
		if len(view.Items) != 1 || !view.Items[0].Available || view.Items[0].Subtotal != 3000 {
			t.Fatalf("unexpected view items: %+v", view.Items)
		}
		if view.Items[0].DiscountPrice != nil {
			t.Fatalf("expected no discount price, got %v", *view.Items[0].DiscountPrice)
		}
		if view.TotalPrice != 3000 {
			t.Fatalf("expected total price 3000 (from original price), got %v", view.TotalPrice)
		}
	})

	t.Run("uses discount price for subtotal and total when available", func(t *testing.T) {
		discountRepo := newFakeProductRepository()
		discountRepo.products["product-2"] = mustNewProductWithDiscount(t, "product-2", "Discounted product", 1000, 800, "UZS", 10)
		discountRepo.products["product-1"] = mustNewProduct(t, "product-1", "Test product", 1000, "UZS", 10)

		discountCartRepo := newFakeCartRepository()
		c := domain.NewCart("user-3")
		_ = c.AddItem("product-1", 2)  // no discount: 2 * 1000 = 2000
		_ = c.AddItem("product-2", 3)  // discount:    3 * 800  = 2400
		_ = discountCartRepo.Save(context.Background(), c)

		discountUC := application.NewGetCartUseCase(discountCartRepo, discountRepo)
		view, err := discountUC.Execute(context.Background(), "user-3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var discountedItem, regularItem application.CartItemView
		for _, item := range view.Items {
			if item.ProductID == "product-2" {
				discountedItem = item
			} else {
				regularItem = item
			}
		}

		if discountedItem.DiscountPrice == nil || *discountedItem.DiscountPrice != 800 {
			t.Fatalf("expected discount price 800, got %+v", discountedItem)
		}
		if discountedItem.UnitPrice != 1000 {
			t.Fatalf("expected original unit price 1000 to still be shown, got %v", discountedItem.UnitPrice)
		}
		if discountedItem.Subtotal != 2400 {
			t.Fatalf("expected subtotal computed from discount price (2400), got %v", discountedItem.Subtotal)
		}
		if regularItem.Subtotal != 2000 {
			t.Fatalf("expected subtotal computed from original price (2000), got %v", regularItem.Subtotal)
		}
		if view.TotalPrice != 4400 {
			t.Fatalf("expected total price 4400 (2000 + 2400), got %v", view.TotalPrice)
		}
	})

	t.Run("deleted product marked unavailable", func(t *testing.T) {
		c := domain.NewCart("user-2")
		_ = c.AddItem("missing-product", 1)
		_ = cartRepo.Save(context.Background(), c)

		view, err := uc.Execute(context.Background(), "user-2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(view.Items) != 1 || view.Items[0].Available {
			t.Fatalf("expected unavailable item, got %+v", view.Items)
		}
	})
}
