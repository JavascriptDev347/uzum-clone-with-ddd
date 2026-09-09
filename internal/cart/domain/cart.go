package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuantity   = errors.New("quantity must be greater than zero")
	ErrItemNotFound      = errors.New("cart item not found")
	ErrEmptyProductID    = errors.New("product ID cannot be empty")
	ErrEmptyUserID       = errors.New("user ID cannot be empty")
	ErrInsufficientStock = errors.New("requested quantity exceeds available stock")
)

type CartItem struct {
	productID string
	quantity  int
}

func NewCartItem(productID string, quantity int) (CartItem, error) {
	if productID == "" {
		return CartItem{}, ErrEmptyProductID
	}
	if quantity <= 0 {
		return CartItem{}, ErrInvalidQuantity
	}
	return CartItem{productID: productID, quantity: quantity}, nil
}

func (i CartItem) ProductID() string { return i.productID }
func (i CartItem) Quantity() int     { return i.quantity }

type Cart struct {
	id        uuid.UUID
	userID    string
	items     []CartItem
	createdAt time.Time
	updatedAt time.Time
}

// NewCart — yangi cart yaratish uchun (use case ichida lazy-create bo'lganda chaqiriladi)
func NewCart(userID string) *Cart {

	now := time.Now()
	return &Cart{
		id:        uuid.New(),
		userID:    userID,
		items:     []CartItem{},
		createdAt: now,
		updatedAt: now,
	}
}

// NewCartFromRepository — DB'dan qayta tiklash uchun, domen konstruktoridan ajratilgan
func NewCartFromRepository(id uuid.UUID, userID string, items []CartItem, createdAt, updatedAt time.Time) *Cart {
	return &Cart{id: id, userID: userID, items: items, createdAt: createdAt, updatedAt: updatedAt}
}

func (c *Cart) ID() uuid.UUID        { return c.id }
func (c *Cart) UserID() string       { return c.userID }
func (c *Cart) Items() []CartItem    { return c.items }
func (c *Cart) CreatedAt() time.Time { return c.createdAt }
func (c *Cart) UpdatedAt() time.Time { return c.updatedAt }

// QuantityOf — berilgan mahsulotning cart'dagi joriy miqdori (mavjud bo'lmasa 0)
func (c *Cart) QuantityOf(productID string) int {
	for i := range c.items {
		if c.items[i].productID == productID {
			return c.items[i].quantity
		}
	}
	return 0
}

// AddItem — upsert: productID mavjud bo'lsa quantity qo'shiladi, aks holda yangi item qo'shiladi
func (c *Cart) AddItem(productID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	for i := range c.items {
		if c.items[i].productID == productID {
			c.items[i].quantity += quantity
			c.updatedAt = time.Now()
			return nil
		}
	}
	item, err := NewCartItem(productID, quantity)
	if err != nil {
		return err
	}
	c.items = append(c.items, item)
	c.updatedAt = time.Now()
	return nil
}

// UpdateItemQuantity — mavjud item'ning quantity'sini to'g'ridan-to'g'ri (absolyut) belgilaydi
func (c *Cart) UpdateItemQuantity(productID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	for i := range c.items {
		if c.items[i].productID == productID {
			c.items[i].quantity = quantity
			c.updatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

// RemoveItem — item'ni butunlay olib tashlaydi
func (c *Cart) RemoveItem(productID string) error {
	for i := range c.items {
		if c.items[i].productID == productID {
			c.items = append(c.items[:i], c.items[i+1:]...)
			c.updatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

// Clear — checkout'dan keyin cart'ni bo'shatish uchun
func (c *Cart) Clear() {
	c.items = []CartItem{}
	c.updatedAt = time.Now()
}

// TotalItems — barcha itemlar quantity'sining yig'indisi (UI badge uchun)
func (c *Cart) TotalItems() int {
	total := 0
	for _, item := range c.items {
		total += item.Quantity()
	}
	return total
}
