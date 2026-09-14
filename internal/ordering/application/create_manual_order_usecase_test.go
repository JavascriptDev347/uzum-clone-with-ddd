package application_test

import (
	"context"
	"errors"
	"testing"

	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

func TestCreateManualOrderUseCase_Execute_HappyPath(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)
	productRepo.products["p2"] = mustNewProduct(t, "p2", "Product 2", 500, "UZS", 10)

	stockReserver := &fakeStockReserver{}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCreateManualOrderUseCase(orderRepo, productRepo, stockReserver)

	order, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
		Items: []application.ManualOrderItemInput{
			{ProductID: "p1", Quantity: 2},
			{ProductID: "p2", Quantity: 3},
		},
		PaymentStatus: domain.PaymentStatusPaid,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.UserID() != nil {
		t.Fatalf("expected nil UserID for a manually created (offline) order, got %v", *order.UserID())
	}
	if order.PaymentStatus() != domain.PaymentStatusPaid {
		t.Fatalf("expected admin-specified PaymentStatus Paid, got %v", order.PaymentStatus())
	}
	if order.DeliveryStatus() != domain.DeliveryStatusPreparing {
		t.Fatalf("expected default DeliveryStatus Preparing, got %v", order.DeliveryStatus())
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

	if _, ok := orderRepo.orders[order.ID().String()]; !ok {
		t.Fatal("expected order to be saved in the repository")
	}
}

func TestCreateManualOrderUseCase_Execute_DefaultsToUnpaidWhenRequested(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)

	uc := application.NewCreateManualOrderUseCase(newFakeOrderRepository(), productRepo, &fakeStockReserver{})

	order, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address:       "Tashkent, Chilonzor 1",
		Phone:         validPhone,
		Items:         []application.ManualOrderItemInput{{ProductID: "p1", Quantity: 1}},
		PaymentStatus: domain.PaymentStatusUnpaid,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.PaymentStatus() != domain.PaymentStatusUnpaid {
		t.Fatalf("expected Unpaid, got %v", order.PaymentStatus())
	}
}

func TestCreateManualOrderUseCase_Execute_EmptyItemsRejected(t *testing.T) {
	uc := application.NewCreateManualOrderUseCase(newFakeOrderRepository(), newFakeProductRepository(), &fakeStockReserver{})

	_, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address:       "Tashkent, Chilonzor 1",
		Phone:         validPhone,
		Items:         nil,
		PaymentStatus: domain.PaymentStatusUnpaid,
	})
	if !errors.Is(err, domain.ErrEmptyOrderItems) {
		t.Fatalf("expected ErrEmptyOrderItems, got %v", err)
	}
}

func TestCreateManualOrderUseCase_Execute_InvalidPaymentStatusRejected(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)
	stockReserver := &fakeStockReserver{}

	uc := application.NewCreateManualOrderUseCase(newFakeOrderRepository(), productRepo, stockReserver)

	_, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address:       "Tashkent, Chilonzor 1",
		Phone:         validPhone,
		Items:         []application.ManualOrderItemInput{{ProductID: "p1", Quantity: 1}},
		PaymentStatus: domain.PaymentStatus("refunded"),
	})
	if !errors.Is(err, domain.ErrInvalidPaymentStatus) {
		t.Fatalf("expected ErrInvalidPaymentStatus, got %v", err)
	}
	if len(stockReserver.decrementCalls) != 0 {
		t.Fatalf("expected no stock decrement before payment status is validated, got %d", len(stockReserver.decrementCalls))
	}
}

func TestCreateManualOrderUseCase_Execute_InsufficientStockRollsBackEarlierDecrements(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)
	productRepo.products["p2"] = mustNewProduct(t, "p2", "Product 2", 500, "UZS", 1)

	stockReserver := &fakeStockReserver{failOnProductID: "p2"}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCreateManualOrderUseCase(orderRepo, productRepo, stockReserver)

	_, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
		Items: []application.ManualOrderItemInput{
			{ProductID: "p1", Quantity: 2},
			{ProductID: "p2", Quantity: 5}, // fails
		},
		PaymentStatus: domain.PaymentStatusPaid,
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
	if len(orderRepo.orders) != 0 {
		t.Fatal("expected no order to be saved on failed creation")
	}
}

func TestCreateManualOrderUseCase_Execute_MergesDuplicateProductIDs(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.products["p1"] = mustNewProduct(t, "p1", "Product 1", 1000, "UZS", 10)
	productRepo.products["p2"] = mustNewProduct(t, "p2", "Product 2", 500, "UZS", 10)

	stockReserver := &fakeStockReserver{}
	orderRepo := newFakeOrderRepository()

	uc := application.NewCreateManualOrderUseCase(orderRepo, productRepo, stockReserver)

	order, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address: "Tashkent, Chilonzor 1",
		Phone:   validPhone,
		Items: []application.ManualOrderItemInput{
			{ProductID: "p1", Quantity: 2},
			{ProductID: "p2", Quantity: 1},
			{ProductID: "p1", Quantity: 3}, // duplicate of the first line - should be merged
		},
		PaymentStatus: domain.PaymentStatusPaid,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order.Items()) != 2 {
		t.Fatalf("expected duplicate ProductID entries to be merged into 2 order items, got %d", len(order.Items()))
	}

	var p1Quantity, p2Quantity int
	for _, item := range order.Items() {
		switch item.ProductID() {
		case "p1":
			p1Quantity = item.Quantity()
		case "p2":
			p2Quantity = item.Quantity()
		}
	}
	if p1Quantity != 5 {
		t.Fatalf("expected merged p1 quantity to be 2+3=5, got %d", p1Quantity)
	}
	if p2Quantity != 1 {
		t.Fatalf("expected p2 quantity to remain 1, got %d", p2Quantity)
	}

	if len(stockReserver.decrementCalls) != 2 {
		t.Fatalf("expected 2 decrement calls (one per unique product), got %d", len(stockReserver.decrementCalls))
	}
	for _, call := range stockReserver.decrementCalls {
		if call.productID == "p1" && call.quantity != 5 {
			t.Fatalf("expected single merged decrement of p1 x5, got %+v", call)
		}
	}
}

func TestCreateManualOrderUseCase_Execute_UnknownProductRejected(t *testing.T) {
	uc := application.NewCreateManualOrderUseCase(newFakeOrderRepository(), newFakeProductRepository(), &fakeStockReserver{})

	_, err := uc.Execute(context.Background(), application.CreateManualOrderInput{
		Address:       "Tashkent, Chilonzor 1",
		Phone:         validPhone,
		Items:         []application.ManualOrderItemInput{{ProductID: "missing", Quantity: 1}},
		PaymentStatus: domain.PaymentStatusUnpaid,
	})
	if !errors.Is(err, catalog.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}
