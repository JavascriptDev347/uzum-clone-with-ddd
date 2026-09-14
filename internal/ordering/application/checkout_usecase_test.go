package application_test

import (
	"context"
	"errors"
	"testing"

	cartapp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/application"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/money"
)

// fakeOrderRepository - domain.OrderRepository ning xotiradagi implementatsiyasi, testlar uchun.
type fakeOrderRepository struct {
	orders  map[string]*domain.Order
	saveErr error
}

func newFakeOrderRepository() *fakeOrderRepository {
	return &fakeOrderRepository{orders: make(map[string]*domain.Order)}
}

func (f *fakeOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.orders[order.ID().String()] = order
	return nil
}

func (f *fakeOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	o, ok := f.orders[id]
	if !ok {
		return nil, nil
	}
	return o, nil
}

func (f *fakeOrderRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	result := make([]*domain.Order, 0)
	for _, o := range f.orders {
		if o.UserID() != nil && *o.UserID() == userID {
			result = append(result, o)
		}
	}
	return result, nil
}

func (f *fakeOrderRepository) FindAll(ctx context.Context) ([]*domain.Order, error) {
	result := make([]*domain.Order, 0, len(f.orders))
	for _, o := range f.orders {
		result = append(result, o)
	}
	return result, nil
}

func (f *fakeOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	f.orders[order.ID().String()] = order
	return nil
}

// fakeCartReader - application.CartReader ning test uchun soxta implementatsiyasi.
type fakeCartReader struct {
	view  *cartapp.CartView
	err   error
	calls int
}

func (f *fakeCartReader) Execute(ctx context.Context, userID string) (*cartapp.CartView, error) {
	f.calls++
	return f.view, f.err
}

// fakeCartClearer - application.CartClearer ning test uchun soxta implementatsiyasi.
type fakeCartClearer struct {
	called       bool
	calledUserID string
	err          error
}

func (f *fakeCartClearer) Execute(ctx context.Context, userID string) error {
	f.called = true
	f.calledUserID = userID
	return f.err
}

type stockCall struct {
	productID string
	quantity  int
}

// fakeStockReserver - domain.StockReserver ning test uchun soxta implementatsiyasi.
type fakeStockReserver struct {
	decrementCalls  []stockCall
	restoreCalls    []stockCall
	failOnProductID string
	failErr         error
}

func (f *fakeStockReserver) DecrementStock(ctx context.Context, productID string, quantity int) error {
	f.decrementCalls = append(f.decrementCalls, stockCall{productID, quantity})
	if productID == f.failOnProductID {
		if f.failErr != nil {
			return f.failErr
		}
		return domain.ErrInsufficientStock
	}
	return nil
}

func (f *fakeStockReserver) RestoreStock(ctx context.Context, productID string, quantity int) error {
	f.restoreCalls = append(f.restoreCalls, stockCall{productID, quantity})
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
		return nil, nil
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

func (f *fakeProductRepository) DecrementStock(ctx context.Context, id string, quantity int) error {
	return nil
}

func (f *fakeProductRepository) IncrementStock(ctx context.Context, id string, quantity int) error {
	return nil
}

func mustNewProduct(t *testing.T, id, nameUz string, amount float64, currency string, stock int) *catalog.Product {
	t.Helper()
	price, err := money.NewMoney(amount, currency)
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

const validPhone = "+998901234567"

func TestCheckoutUseCase_Execute_HappyPath(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)
	productRepo.products["p2"] = mustNewProduct(t, "p2", "Product 2", 500, "UZS", 10)

	cartReader := &fakeCartReader{view: &cartapp.CartView{
		Items: []cartapp.CartItemView{
			{ProductID: "p1", Quantity: 2, Available: true},
			{ProductID: "p2", Quantity: 3, Available: true},
		},
		TotalItems: 5,
	}}
	cartClearer := &fakeCartClearer{}
	stockReserver := &fakeStockReserver{}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCheckoutUseCase(orderRepo, cartReader, cartClearer, productRepo, stockReserver)

	order, err := uc.Execute(context.Background(), application.CheckoutInput{
		UserID:  "user-1",
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order.Items()) != 2 {
		t.Fatalf("expected 2 order items, got %d", len(order.Items()))
	}

	total, err := order.Total()
	if err != nil {
		t.Fatalf("unexpected error computing total: %v", err)
	}
	if total.Amount() != 3500 { // 2*1000 + 3*500
		t.Fatalf("expected total 3500, got %v", total.Amount())
	}

	if len(stockReserver.decrementCalls) != 2 {
		t.Fatalf("expected 2 decrement calls, got %d", len(stockReserver.decrementCalls))
	}
	if len(stockReserver.restoreCalls) != 0 {
		t.Fatalf("expected no restore calls on happy path, got %d", len(stockReserver.restoreCalls))
	}

	if !cartClearer.called || cartClearer.calledUserID != "user-1" {
		t.Fatalf("expected cart to be cleared for user-1, got called=%v userID=%v", cartClearer.called, cartClearer.calledUserID)
	}

	if _, ok := orderRepo.orders[order.ID().String()]; !ok {
		t.Fatal("expected order to be saved in the repository")
	}
}

func TestCheckoutUseCase_Execute_EmptyCartRejected(t *testing.T) {
	productRepo := newFakeProductRepository()
	cartReader := &fakeCartReader{view: &cartapp.CartView{Items: []cartapp.CartItemView{}}}
	cartClearer := &fakeCartClearer{}
	stockReserver := &fakeStockReserver{}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCheckoutUseCase(orderRepo, cartReader, cartClearer, productRepo, stockReserver)

	_, err := uc.Execute(context.Background(), application.CheckoutInput{
		UserID:  "user-1",
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
	})
	if !errors.Is(err, domain.ErrEmptyOrderItems) {
		t.Fatalf("expected ErrEmptyOrderItems, got %v", err)
	}
	if len(stockReserver.decrementCalls) != 0 {
		t.Fatalf("expected no stock decrements for an empty cart, got %d", len(stockReserver.decrementCalls))
	}
	if cartClearer.called {
		t.Fatal("expected cart clearer not to be called")
	}
	if len(orderRepo.orders) != 0 {
		t.Fatal("expected no order to be saved")
	}
}

func TestCheckoutUseCase_Execute_InvalidCustomerInfoShortCircuits(t *testing.T) {
	productRepo := newFakeProductRepository()
	cartReader := &fakeCartReader{view: &cartapp.CartView{Items: []cartapp.CartItemView{{ProductID: "p1", Quantity: 1, Available: true}}}}
	cartClearer := &fakeCartClearer{}
	stockReserver := &fakeStockReserver{}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCheckoutUseCase(orderRepo, cartReader, cartClearer, productRepo, stockReserver)

	_, err := uc.Execute(context.Background(), application.CheckoutInput{
		UserID:  "user-1",
		Address: "Tashkent, Chilonzor 1",
		Phone:   "not-a-phone",
	})
	if !errors.Is(err, domain.ErrInvalidPhone) {
		t.Fatalf("expected ErrInvalidPhone, got %v", err)
	}
	if cartReader.calls != 0 {
		t.Fatalf("expected cart reader not to be called before customer info is validated, got %d calls", cartReader.calls)
	}
}

func TestCheckoutUseCase_Execute_UnavailableItemRejected(t *testing.T) {
	productRepo := newFakeProductRepository()
	cartReader := &fakeCartReader{view: &cartapp.CartView{
		Items: []cartapp.CartItemView{{ProductID: "missing", Quantity: 1, Available: false}},
	}}
	cartClearer := &fakeCartClearer{}
	stockReserver := &fakeStockReserver{}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCheckoutUseCase(orderRepo, cartReader, cartClearer, productRepo, stockReserver)

	_, err := uc.Execute(context.Background(), application.CheckoutInput{
		UserID:  "user-1",
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
	})
	if !errors.Is(err, catalog.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestCheckoutUseCase_Execute_InsufficientStockRollsBackEarlierDecrements(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)
	productRepo.products["p2"] = mustNewProduct(t, "p2", "Product 2", 500, "UZS", 1)

	cartReader := &fakeCartReader{view: &cartapp.CartView{
		Items: []cartapp.CartItemView{
			{ProductID: "p1", Quantity: 2, Available: true},
			{ProductID: "p2", Quantity: 5, Available: true}, // this one fails
		},
	}}
	cartClearer := &fakeCartClearer{}
	stockReserver := &fakeStockReserver{failOnProductID: "p2"}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCheckoutUseCase(orderRepo, cartReader, cartClearer, productRepo, stockReserver)

	_, err := uc.Execute(context.Background(), application.CheckoutInput{
		UserID:  "user-1",
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
	})
	if !errors.Is(err, domain.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}

	if len(stockReserver.restoreCalls) != 1 {
		t.Fatalf("expected exactly 1 rollback (for p1), got %d", len(stockReserver.restoreCalls))
	}
	if stockReserver.restoreCalls[0].productID != "p1" || stockReserver.restoreCalls[0].quantity != 2 {
		t.Fatalf("expected rollback of p1 x2, got %+v", stockReserver.restoreCalls[0])
	}

	if cartClearer.called {
		t.Fatal("expected cart not to be cleared on failed checkout")
	}
	if len(orderRepo.orders) != 0 {
		t.Fatal("expected no order to be saved on failed checkout")
	}
}

func TestCheckoutUseCase_Execute_UsesDiscountPriceWhenPresent(t *testing.T) {
	price, _ := money.NewMoney(1000, "UZS")
	discount, _ := money.NewMoney(700, "UZS")
	product, err := catalog.NewProduct(catalog.NewProductParams{
		ID:            "p1",
		NameUz:        "Product 1",
		NameEng:       "Product 1",
		NameRu:        "Product 1",
		CategoryID:    "category-1",
		Price:         price,
		DiscountPrice: &discount,
		Slug:          "product-1",
		Stock:         10,
	})
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}

	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = product

	cartReader := &fakeCartReader{view: &cartapp.CartView{
		Items: []cartapp.CartItemView{{ProductID: "p1", Quantity: 2, Available: true}},
	}}
	uc := application.NewCheckoutUseCase(newFakeOrderRepository(), cartReader, &fakeCartClearer{}, productRepo, &fakeStockReserver{})

	order, err := uc.Execute(context.Background(), application.CheckoutInput{
		UserID:  "user-1",
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.Items()[0].UnitPrice().Amount() != 700 {
		t.Fatalf("expected snapshotted unit price to use the discount (700), got %v", order.Items()[0].UnitPrice().Amount())
	}
}

func TestGetUserOrdersUseCase_Execute(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	info, _ := domain.NewCustomerInfo("Tashkent, Chilonzor 1", validPhone, "")
	item, _ := domain.NewOrderItem("p1", "Product 1", func() money.Money { m, _ := money.NewMoney(1000, "UZS"); return m }(), 1)

	userA := "user-a"
	orderA, _ := domain.NewOrder(&userA, info, []domain.OrderItem{item})
	_ = orderRepo.Save(context.Background(), orderA)

	userB := "user-b"
	orderB, _ := domain.NewOrder(&userB, info, []domain.OrderItem{item})
	_ = orderRepo.Save(context.Background(), orderB)

	uc := application.NewGetUserOrdersUseCase(orderRepo)

	orders, err := uc.Execute(context.Background(), "user-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 1 || orders[0].ID() != orderA.ID() {
		t.Fatalf("expected only user-a's order, got %+v", orders)
	}
}

func TestGetOrderByIDUseCase_Execute(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	info, _ := domain.NewCustomerInfo("Tashkent, Chilonzor 1", validPhone, "")
	item, _ := domain.NewOrderItem("p1", "Product 1", func() money.Money { m, _ := money.NewMoney(1000, "UZS"); return m }(), 1)

	owner := "owner-1"
	order, _ := domain.NewOrder(&owner, info, []domain.OrderItem{item})
	_ = orderRepo.Save(context.Background(), order)

	uc := application.NewGetOrderByIDUseCase(orderRepo)

	t.Run("owner can access their own order", func(t *testing.T) {
		got, err := uc.Execute(context.Background(), order.ID().String(), "owner-1", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID() != order.ID() {
			t.Fatalf("expected owner's order, got %+v", got)
		}
	})

	t.Run("non-owner is denied", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), order.ID().String(), "someone-else", false)
		if !errors.Is(err, domain.ErrOrderAccessDenied) {
			t.Fatalf("expected ErrOrderAccessDenied, got %v", err)
		}
	})

	t.Run("admin can access any order", func(t *testing.T) {
		got, err := uc.Execute(context.Background(), order.ID().String(), "someone-else", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID() != order.ID() {
			t.Fatalf("expected order to be returned for admin, got %+v", got)
		}
	})

	t.Run("unknown id returns ErrOrderNotFound", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "does-not-exist", "owner-1", false)
		if !errors.Is(err, domain.ErrOrderNotFound) {
			t.Fatalf("expected ErrOrderNotFound, got %v", err)
		}
	})

	t.Run("offline order (nil UserID) is denied to any non-admin", func(t *testing.T) {
		offlineOrder, _ := domain.NewOrder(nil, info, []domain.OrderItem{item})
		_ = orderRepo.Save(context.Background(), offlineOrder)

		_, err := uc.Execute(context.Background(), offlineOrder.ID().String(), "owner-1", false)
		if !errors.Is(err, domain.ErrOrderAccessDenied) {
			t.Fatalf("expected ErrOrderAccessDenied, got %v", err)
		}
	})
}
