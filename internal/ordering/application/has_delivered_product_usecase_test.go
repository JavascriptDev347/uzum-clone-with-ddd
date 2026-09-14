package application_test

import (
	"context"
	"testing"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

func TestHasDeliveredProductUseCase_Execute(t *testing.T) {
	orderRepo := newFakeOrderRepository()

	info, err := domain.NewCustomerInfo("Tashkent, Chilonzor 1", validPhone, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	userID := "user-1"

	deliveredItem := mustOrderItemForTest(t, "p1", 1000, "UZS", 1)
	deliveredOrder, err := domain.NewOrder(&userID, info, []domain.OrderItem{deliveredItem})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := deliveredOrder.AdvanceDelivery(domain.DeliveryStatusDelivered); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := orderRepo.Save(context.Background(), deliveredOrder); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preparingItem := mustOrderItemForTest(t, "p2", 500, "UZS", 1)
	preparingOrder, err := domain.NewOrder(&userID, info, []domain.OrderItem{preparingItem})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := orderRepo.Save(context.Background(), preparingOrder); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	uc := application.NewHasDeliveredProductUseCase(orderRepo)

	t.Run("delivered product is found", func(t *testing.T) {
		orderID, ok, err := uc.Execute(context.Background(), "user-1", "p1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true for a delivered product")
		}
		if orderID != deliveredOrder.ID().String() {
			t.Fatalf("expected orderID %v, got %v", deliveredOrder.ID().String(), orderID)
		}
	})

	t.Run("product only in a non-delivered order is not found", func(t *testing.T) {
		_, ok, err := uc.Execute(context.Background(), "user-1", "p2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for a product whose order is still Preparing")
		}
	})

	t.Run("unknown product is not found", func(t *testing.T) {
		_, ok, err := uc.Execute(context.Background(), "user-1", "p999")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for a product the user never ordered")
		}
	})

	t.Run("unknown user is not found", func(t *testing.T) {
		_, ok, err := uc.Execute(context.Background(), "user-2", "p1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for a different user")
		}
	})
}

func mustOrderItemForTest(t *testing.T, productID string, amount float64, currency string, quantity int) domain.OrderItem {
	t.Helper()
	item, err := domain.NewOrderItem(productID, "Product "+productID, mustMoney(t, amount, currency), quantity)
	if err != nil {
		t.Fatalf("unexpected error building OrderItem: %v", err)
	}
	return item
}
