package domain

import (
	"errors"
	"testing"
)

func TestNewCartItem(t *testing.T) {
	t.Run("valid item", func(t *testing.T) {
		item, err := NewCartItem("product-1", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ProductID() != "product-1" || item.Quantity() != 2 {
			t.Fatalf("unexpected item: %+v", item)
		}
	})

	t.Run("empty product id", func(t *testing.T) {
		if _, err := NewCartItem("", 1); !errors.Is(err, ErrEmptyProductID) {
			t.Fatalf("expected ErrEmptyProductID, got %v", err)
		}
	})

	t.Run("invalid quantity", func(t *testing.T) {
		if _, err := NewCartItem("product-1", 0); !errors.Is(err, ErrInvalidQuantity) {
			t.Fatalf("expected ErrInvalidQuantity, got %v", err)
		}
	})
}

func TestCart_AddItem(t *testing.T) {
	c := NewCart("user-1")

	if err := c.AddItem("product-1", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.TotalItems() != 2 {
		t.Fatalf("expected total 2, got %d", c.TotalItems())
	}

	// adding the same product again should increment quantity, not duplicate
	if err := c.AddItem("product-1", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Items()) != 1 {
		t.Fatalf("expected 1 distinct item, got %d", len(c.Items()))
	}
	if c.TotalItems() != 5 {
		t.Fatalf("expected total 5, got %d", c.TotalItems())
	}

	if err := c.AddItem("product-2", 0); !errors.Is(err, ErrInvalidQuantity) {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestCart_UpdateItemQuantity(t *testing.T) {
	c := NewCart("user-1")
	_ = c.AddItem("product-1", 2)

	if err := c.UpdateItemQuantity("product-1", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.TotalItems() != 10 {
		t.Fatalf("expected total 10, got %d", c.TotalItems())
	}

	if err := c.UpdateItemQuantity("product-1", 0); !errors.Is(err, ErrInvalidQuantity) {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}

	if err := c.UpdateItemQuantity("missing", 1); !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}

func TestCart_RemoveItem(t *testing.T) {
	c := NewCart("user-1")
	_ = c.AddItem("product-1", 2)
	_ = c.AddItem("product-2", 1)

	if err := c.RemoveItem("product-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Items()) != 1 {
		t.Fatalf("expected 1 item left, got %d", len(c.Items()))
	}

	if err := c.RemoveItem("product-1"); !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}

func TestCart_Clear(t *testing.T) {
	c := NewCart("user-1")
	_ = c.AddItem("product-1", 2)

	c.Clear()

	if len(c.Items()) != 0 || c.TotalItems() != 0 {
		t.Fatalf("expected empty cart after clear, got items=%v total=%d", c.Items(), c.TotalItems())
	}
}
