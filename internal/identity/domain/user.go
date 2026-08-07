package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyUserID       = errors.New("identity: user id bo'sh bo'lishi mumkin emas")
	ErrEmptyPasswordHash = errors.New("identity: password hash bo'sh bo'lishi mumkin emas")
)

// User - Identity context'ining Entity'i. Identity - ID orqali aniqlanadi.
type User struct {
	id           string
	email        Email
	passwordHash string
	role         Role
	createdAt    time.Time
}

// NewUser - yangi User yaratadi (Register jarayonida chaqiriladi).
func NewUser(id string, email Email, passwordHash string, role Role, createdAt time.Time) (*User, error) {
	if id == "" {
		return nil, ErrEmptyUserID
	}

	if passwordHash == "" {
		return nil, ErrEmptyPasswordHash
	}
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}
	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
		createdAt:    time.Now(),
	}, nil
}

// NewUserFromRepository - DB qatoridan User qayta tiklash uchun.
// Validatsiya YO'Q - ma'lumot allaqachon bir marta tekshirilgan.
func NewUserFromRepository(id string, email Email, passwordHash string, role Role, createdAt time.Time) *User {
	return &User{id: id, email: email, passwordHash: passwordHash, role: role, createdAt: createdAt}
}
func (u *User) ID() string           { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) Role() Role           { return u.role }
func (u *User) CreatedAt() time.Time { return u.createdAt }

// ChangeRole - user role'ini o'zgartiradi.
func (u *User) ChangeRole(newRole Role) error {
	if !newRole.IsValid() {
		return ErrInvalidRole
	}
	u.role = newRole
	return nil
}

// ChangePasswordHash - user password hash'ini o'zgartiradi.
func (u *User) ChangePasswordHash(newHash string) error {
	if newHash == "" {
		return ErrEmptyPasswordHash
	}
	u.passwordHash = newHash
	return nil
}
