package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/money"
)

func mustOrderForAdminTests(t *testing.T, repo *fakeOrderRepository, userID string) *domain.Order {
	t.Helper()
	info, err := domain.NewCustomerInfo("Tashkent, Chilonzor 1", validPhone, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	price, err := money.NewMoney(1000, "UZS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item, err := domain.NewOrderItem("p1", "Product 1", price, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	order, err := domain.NewOrder(&userID, info, []domain.OrderItem{item})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.Save(context.Background(), order); err != nil {
		t.Fatalf("unexpected error saving order: %v", err)
	}
	return order
}

func TestUpdatePaymentStatusUseCase_Execute(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	order := mustOrderForAdminTests(t, orderRepo, "user-1")
	uc := application.NewUpdatePaymentStatusUseCase(orderRepo)

	t.Run("marks paid", func(t *testing.T) {
		updated, err := uc.Execute(context.Background(), order.ID().String(), domain.PaymentStatusPaid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.PaymentStatus() != domain.PaymentStatusPaid {
			t.Fatalf("expected Paid, got %v", updated.PaymentStatus())
		}
		persisted, _ := orderRepo.FindByID(context.Background(), order.ID().String())
		if persisted.PaymentStatus() != domain.PaymentStatusPaid {
			t.Fatalf("expected persisted status Paid, got %v", persisted.PaymentStatus())
		}
	})

	t.Run("marks unpaid", func(t *testing.T) {
		updated, err := uc.Execute(context.Background(), order.ID().String(), domain.PaymentStatusUnpaid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.PaymentStatus() != domain.PaymentStatusUnpaid {
			t.Fatalf("expected Unpaid, got %v", updated.PaymentStatus())
		}
	})

	t.Run("invalid status rejected", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), order.ID().String(), domain.PaymentStatus("refunded"))
		if !errors.Is(err, domain.ErrInvalidPaymentStatus) {
			t.Fatalf("expected ErrInvalidPaymentStatus, got %v", err)
		}
	})

	t.Run("unknown order returns ErrOrderNotFound", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "does-not-exist", domain.PaymentStatusPaid)
		if !errors.Is(err, domain.ErrOrderNotFound) {
			t.Fatalf("expected ErrOrderNotFound, got %v", err)
		}
	})
}

func TestUpdateDeliveryStatusUseCase_Execute(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	order := mustOrderForAdminTests(t, orderRepo, "user-1")
	uc := application.NewUpdateDeliveryStatusUseCase(orderRepo)

	t.Run("advances forward", func(t *testing.T) {
		updated, err := uc.Execute(context.Background(), order.ID().String(), domain.DeliveryStatusHandedToCourier)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.DeliveryStatus() != domain.DeliveryStatusHandedToCourier {
			t.Fatalf("expected HandedToCourier, got %v", updated.DeliveryStatus())
		}
		persisted, _ := orderRepo.FindByID(context.Background(), order.ID().String())
		if persisted.DeliveryStatus() != domain.DeliveryStatusHandedToCourier {
			t.Fatalf("expected persisted status HandedToCourier, got %v", persisted.DeliveryStatus())
		}
	})

	t.Run("rejects moving backwards", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), order.ID().String(), domain.DeliveryStatusPreparing)
		if !errors.Is(err, domain.ErrInvalidDeliveryTransition) {
			t.Fatalf("expected ErrInvalidDeliveryTransition, got %v", err)
		}
		// status must remain unchanged after the rejected transition
		persisted, _ := orderRepo.FindByID(context.Background(), order.ID().String())
		if persisted.DeliveryStatus() != domain.DeliveryStatusHandedToCourier {
			t.Fatalf("expected status to remain HandedToCourier, got %v", persisted.DeliveryStatus())
		}
	})

	t.Run("invalid status value rejected", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), order.ID().String(), domain.DeliveryStatus("shipped"))
		if !errors.Is(err, domain.ErrInvalidDeliveryStatus) {
			t.Fatalf("expected ErrInvalidDeliveryStatus, got %v", err)
		}
	})

	t.Run("unknown order returns ErrOrderNotFound", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "does-not-exist", domain.DeliveryStatusDelivered)
		if !errors.Is(err, domain.ErrOrderNotFound) {
			t.Fatalf("expected ErrOrderNotFound, got %v", err)
		}
	})
}

func TestGetAllOrdersUseCase_Execute(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	mustOrderForAdminTests(t, orderRepo, "user-1")
	mustOrderForAdminTests(t, orderRepo, "user-2")

	uc := application.NewGetAllOrdersUseCase(orderRepo)

	orders, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}
}
