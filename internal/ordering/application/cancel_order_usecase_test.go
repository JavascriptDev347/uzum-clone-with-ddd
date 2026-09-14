package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/money"
)

func mustMoney(t *testing.T, amount float64, currency string) money.Money {
	t.Helper()
	m, err := money.NewMoney(amount, currency)
	if err != nil {
		t.Fatalf("unexpected error building Money: %v", err)
	}
	return m
}

func mustSeedOrder(t *testing.T, orderRepo *fakeOrderRepository, items []domain.OrderItem, deliveryStatus domain.DeliveryStatus) *domain.Order {
	t.Helper()

	info, err := domain.NewCustomerInfo("Tashkent, Chilonzor 1", validPhone, "")
	if err != nil {
		t.Fatalf("unexpected error building CustomerInfo: %v", err)
	}

	order, err := domain.NewOrder(nil, info, items)
	if err != nil {
		t.Fatalf("unexpected error building order: %v", err)
	}

	if err := orderRepo.Save(context.Background(), order); err != nil {
		t.Fatalf("unexpected error saving order: %v", err)
	}

	// AdvanceDelivery bilan Preparing'dan boshqa holatga o'tkazamiz - testda boshlang'ich
	// holatni to'g'ridan-to'g'ri o'rnatish kerak bo'lganda (masalan HandedToCourier/Delivered).
	if deliveryStatus != domain.DeliveryStatusPreparing {
		if err := order.AdvanceDelivery(deliveryStatus); err != nil {
			t.Fatalf("unexpected error advancing delivery status to %v: %v", deliveryStatus, err)
		}
	}

	return order
}

func TestCancelOrderUseCase_Execute_PreparingOrderIsCancelledAndStockRestored(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	stockReserver := &fakeStockReserver{}

	item1, err := domain.NewOrderItem("p1", "Product 1", mustMoney(t, 1000, "UZS"), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item2, err := domain.NewOrderItem("p2", "Product 2", mustMoney(t, 500, "UZS"), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	order := mustSeedOrder(t, orderRepo, []domain.OrderItem{item1, item2}, domain.DeliveryStatusPreparing)

	uc := application.NewCancelOrderUseCase(orderRepo, stockReserver)

	cancelled, err := uc.Execute(context.Background(), order.ID().String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cancelled.DeliveryStatus() != domain.DeliveryStatusCancelled {
		t.Fatalf("expected DeliveryStatus Cancelled, got %v", cancelled.DeliveryStatus())
	}

	if len(stockReserver.restoreCalls) != 2 {
		t.Fatalf("expected 2 restore calls (one per item), got %d", len(stockReserver.restoreCalls))
	}
	restored := map[string]int{}
	for _, call := range stockReserver.restoreCalls {
		restored[call.productID] = call.quantity
	}
	if restored["p1"] != 2 {
		t.Fatalf("expected p1 restored by exactly 2, got %d", restored["p1"])
	}
	if restored["p2"] != 3 {
		t.Fatalf("expected p2 restored by exactly 3, got %d", restored["p2"])
	}

	persisted, ok := orderRepo.orders[order.ID().String()]
	if !ok {
		t.Fatal("expected order to remain in the repository")
	}
	if persisted.DeliveryStatus() != domain.DeliveryStatusCancelled {
		t.Fatalf("expected persisted order to be Cancelled, got %v", persisted.DeliveryStatus())
	}
}

func TestCancelOrderUseCase_Execute_ShippedOrdersAreRejectedAndStockUntouched(t *testing.T) {
	tests := []struct {
		name           string
		deliveryStatus domain.DeliveryStatus
	}{
		{name: "handed to courier", deliveryStatus: domain.DeliveryStatusHandedToCourier},
		{name: "delivered", deliveryStatus: domain.DeliveryStatusDelivered},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newFakeOrderRepository()
			stockReserver := &fakeStockReserver{}

			item, err := domain.NewOrderItem("p1", "Product 1", mustMoney(t, 1000, "UZS"), 2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			order := mustSeedOrder(t, orderRepo, []domain.OrderItem{item}, tt.deliveryStatus)

			uc := application.NewCancelOrderUseCase(orderRepo, stockReserver)

			_, err = uc.Execute(context.Background(), order.ID().String())
			if !errors.Is(err, domain.ErrOrderAlreadyShipped) {
				t.Fatalf("expected ErrOrderAlreadyShipped, got %v", err)
			}

			if len(stockReserver.restoreCalls) != 0 {
				t.Fatalf("expected no stock restore calls, got %d", len(stockReserver.restoreCalls))
			}

			persisted := orderRepo.orders[order.ID().String()]
			if persisted.DeliveryStatus() != tt.deliveryStatus {
				t.Fatalf("expected order status to remain %v, got %v", tt.deliveryStatus, persisted.DeliveryStatus())
			}
		})
	}
}

func TestCancelOrderUseCase_Execute_UnknownOrderRejected(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	stockReserver := &fakeStockReserver{}

	uc := application.NewCancelOrderUseCase(orderRepo, stockReserver)

	_, err := uc.Execute(context.Background(), "missing-order-id")
	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
	if len(stockReserver.restoreCalls) != 0 {
		t.Fatalf("expected no stock restore calls, got %d", len(stockReserver.restoreCalls))
	}
}
