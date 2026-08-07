package domain

import "errors"

var ErrInvalidRole = errors.New("identity: role noto'g'ri (customer yoki admin bo'lishi kerak)")

type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

// IsValid - Role ruxsat etilgan qiymatlardan biri ekanligini tekshiradi.
func (r Role) IsValid() bool {
	return r == RoleCustomer || r == RoleAdmin
}
