package domain

import (
	"context"
	"errors"
)

type UserRepository interface {
	// Save - yangi User yaratadi
	Save(ctx context.Context, user *User) error

	// FindByID - ID bo'yicha topadi. Topilmasa ErrUserNotFound qaytaradi.
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail - Email bo'yicha topadi (login uchun kerak). Topilmasa ErrUserNotFound.
	FindByEmail(ctx context.Context, email Email) (*User, error)
}

// ErrUserNotFound - domain darajasidagi xato. Postgres'ning sql.ErrNoRows
var ErrUserNotFound = errors.New("user not found")
