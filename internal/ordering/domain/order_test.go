package domain

import (
	"errors"
	"testing"

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

func mustCustomerInfo(t *testing.T) CustomerInfo {
	t.Helper()
	info, err := NewCustomerInfo("Tashkent, Chilonzor 1", "+998901234567", "")
	if err != nil {
		t.Fatalf("unexpected error building CustomerInfo: %v", err)
	}
	return info
}

func mustOrderItem(t *testing.T, productID string, unitPrice money.Money, quantity int) OrderItem {
	t.Helper()
	item, err := NewOrderItem(productID, "Product "+productID, unitPrice, quantity)
	if err != nil {
		t.Fatalf("unexpected error building OrderItem: %v", err)
	}
	return item
}

func TestNewOrderItem(t *testing.T) {
	price := mustMoney(t, 1000, "UZS")

	tests := []struct {
		name        string
		productID   string
		productName string
		quantity    int
		wantErr     error
	}{
		{name: "valid", productID: "p1", productName: "Product 1", quantity: 2, wantErr: nil},
		{name: "empty product id", productID: "", productName: "Product 1", quantity: 2, wantErr: ErrEmptyOrderItemProductID},
		{name: "empty product name", productID: "p1", productName: "", quantity: 2, wantErr: ErrEmptyOrderItemName},
		{name: "zero quantity", productID: "p1", productName: "Product 1", quantity: 0, wantErr: ErrInvalidQuantity},
		{name: "negative quantity", productID: "p1", productName: "Product 1", quantity: -1, wantErr: ErrInvalidQuantity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOrderItem(tt.productID, tt.productName, price, tt.quantity)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewOrder_EmptyItemsRejected(t *testing.T) {
	info := mustCustomerInfo(t)
	validItem := mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1)

	tests := []struct {
		name    string
		items   []OrderItem
		wantErr error
	}{
		{name: "nil items", items: nil, wantErr: ErrEmptyOrderItems},
		{name: "empty slice", items: []OrderItem{}, wantErr: ErrEmptyOrderItems},
		{name: "one item is valid", items: []OrderItem{validItem}, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(nil, info, tt.items)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}
			if order.PaymentStatus() != PaymentStatusUnpaid {
				t.Fatalf("expected default payment status Unpaid, got %v", order.PaymentStatus())
			}
			if order.DeliveryStatus() != DeliveryStatusPreparing {
				t.Fatalf("expected default delivery status Preparing, got %v", order.DeliveryStatus())
			}
		})
	}
}

func TestOrder_UserID_NilForOfflineOrders(t *testing.T) {
	info := mustCustomerInfo(t)
	item := mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1)

	order, err := NewOrder(nil, info, []OrderItem{item})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.UserID() != nil {
		t.Fatalf("expected nil UserID for offline order, got %v", *order.UserID())
	}

	userID := "user-1"
	order2, err := NewOrder(&userID, info, []OrderItem{item})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order2.UserID() == nil || *order2.UserID() != userID {
		t.Fatalf("expected UserID %q, got %v", userID, order2.UserID())
	}
}

func TestOrder_MarkPaid_MarkUnpaid(t *testing.T) {
	info := mustCustomerInfo(t)
	item := mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1)
	order, _ := NewOrder(nil, info, []OrderItem{item})

	order.MarkPaid()
	if order.PaymentStatus() != PaymentStatusPaid {
		t.Fatalf("expected Paid, got %v", order.PaymentStatus())
	}

	order.MarkUnpaid()
	if order.PaymentStatus() != PaymentStatusUnpaid {
		t.Fatalf("expected Unpaid, got %v", order.PaymentStatus())
	}
}

func TestOrder_AdvanceDelivery(t *testing.T) {
	info := mustCustomerInfo(t)
	item := mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1)

	tests := []struct {
		name    string
		from    DeliveryStatus
		to      DeliveryStatus
		wantErr error
	}{
		{name: "preparing to handed to courier", from: DeliveryStatusPreparing, to: DeliveryStatusHandedToCourier, wantErr: nil},
		{name: "preparing to delivered (skip ahead)", from: DeliveryStatusPreparing, to: DeliveryStatusDelivered, wantErr: nil},
		{name: "handed to courier to delivered", from: DeliveryStatusHandedToCourier, to: DeliveryStatusDelivered, wantErr: nil},
		{name: "same status is a no-op", from: DeliveryStatusHandedToCourier, to: DeliveryStatusHandedToCourier, wantErr: nil},
		{name: "handed to courier back to preparing", from: DeliveryStatusHandedToCourier, to: DeliveryStatusPreparing, wantErr: ErrInvalidDeliveryTransition},
		{name: "delivered back to preparing", from: DeliveryStatusDelivered, to: DeliveryStatusPreparing, wantErr: ErrInvalidDeliveryTransition},
		{name: "delivered back to handed to courier", from: DeliveryStatusDelivered, to: DeliveryStatusHandedToCourier, wantErr: ErrInvalidDeliveryTransition},
		{name: "invalid status value", from: DeliveryStatusPreparing, to: DeliveryStatus("shipped"), wantErr: ErrInvalidDeliveryStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(nil, info, []OrderItem{item})
			if err != nil {
				t.Fatalf("unexpected error building order: %v", err)
			}
			order.deliveryStatus = tt.from

			err = order.AdvanceDelivery(tt.to)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr != nil {
				if order.DeliveryStatus() != tt.from {
					t.Fatalf("expected status to remain %v after rejected transition, got %v", tt.from, order.DeliveryStatus())
				}
				return
			}
			if order.DeliveryStatus() != tt.to {
				t.Fatalf("expected status %v, got %v", tt.to, order.DeliveryStatus())
			}
		})
	}
}

func TestOrder_AdvanceDelivery_CancelledIsUnreachableAndTerminal(t *testing.T) {
	info := mustCustomerInfo(t)
	item := mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1)

	tests := []struct {
		name string
		from DeliveryStatus
		to   DeliveryStatus
	}{
		{name: "preparing to cancelled via AdvanceDelivery is rejected", from: DeliveryStatusPreparing, to: DeliveryStatusCancelled},
		{name: "cancelled cannot advance to handed to courier", from: DeliveryStatusCancelled, to: DeliveryStatusHandedToCourier},
		{name: "cancelled cannot advance to delivered", from: DeliveryStatusCancelled, to: DeliveryStatusDelivered},
		{name: "cancelled cannot be re-set to cancelled", from: DeliveryStatusCancelled, to: DeliveryStatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(nil, info, []OrderItem{item})
			if err != nil {
				t.Fatalf("unexpected error building order: %v", err)
			}
			order.deliveryStatus = tt.from

			err = order.AdvanceDelivery(tt.to)
			if !errors.Is(err, ErrInvalidDeliveryTransition) {
				t.Fatalf("expected ErrInvalidDeliveryTransition, got %v", err)
			}
			if order.DeliveryStatus() != tt.from {
				t.Fatalf("expected status to remain %v after rejected transition, got %v", tt.from, order.DeliveryStatus())
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	info := mustCustomerInfo(t)
	item := mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1)

	tests := []struct {
		name    string
		from    DeliveryStatus
		wantErr error
	}{
		{name: "preparing can be cancelled", from: DeliveryStatusPreparing, wantErr: nil},
		{name: "handed to courier cannot be cancelled", from: DeliveryStatusHandedToCourier, wantErr: ErrOrderAlreadyShipped},
		{name: "delivered cannot be cancelled", from: DeliveryStatusDelivered, wantErr: ErrOrderAlreadyShipped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(nil, info, []OrderItem{item})
			if err != nil {
				t.Fatalf("unexpected error building order: %v", err)
			}
			order.deliveryStatus = tt.from

			err = order.Cancel()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				if order.DeliveryStatus() != tt.from {
					t.Fatalf("expected status to remain %v after rejected cancel, got %v", tt.from, order.DeliveryStatus())
				}
				return
			}
			if order.DeliveryStatus() != DeliveryStatusCancelled {
				t.Fatalf("expected status Cancelled, got %v", order.DeliveryStatus())
			}
		})
	}
}

func TestOrder_Total(t *testing.T) {
	info := mustCustomerInfo(t)

	tests := []struct {
		name       string
		items      []OrderItem
		wantAmount float64
		wantErr    error
	}{
		{
			name: "single item",
			items: []OrderItem{
				mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 3),
			},
			wantAmount: 3000,
		},
		{
			name: "multiple items, same currency",
			items: []OrderItem{
				mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 2), // 2000
				mustOrderItem(t, "p2", mustMoney(t, 5000, "UZS"), 1), // 5000
				mustOrderItem(t, "p3", mustMoney(t, 250, "UZS"), 4),  // 1000
			},
			wantAmount: 8000,
		},
		{
			name: "mismatched currencies error out",
			items: []OrderItem{
				mustOrderItem(t, "p1", mustMoney(t, 1000, "UZS"), 1),
				mustOrderItem(t, "p2", mustMoney(t, 10, "USD"), 1),
			},
			wantErr: money.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(nil, info, tt.items)
			if err != nil {
				t.Fatalf("unexpected error building order: %v", err)
			}

			total, err := order.Total()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}
			if total.Amount() != tt.wantAmount {
				t.Fatalf("expected total %v, got %v", tt.wantAmount, total.Amount())
			}
			if total.Currency() != "UZS" {
				t.Fatalf("expected currency UZS, got %v", total.Currency())
			}
		})
	}
}
