package domain

import (
	"regexp"
	"strings"
)

// phonePattern - O'zbekiston telefon raqami formati: +998 bilan boshlanib, keyin 9 ta raqam.
var phonePattern = regexp.MustCompile(`^\+998\d{9}$`)

// CustomerInfo - Value Object. Buyurtmani yetkazib berish uchun kerakli mijoz ma'lumotlari.
// Note ixtiyoriy (bo'sh bo'lishi mumkin), Address va Phone majburiy.
type CustomerInfo struct {
	address string
	phone   string
	note    string
}

func NewCustomerInfo(address, phone, note string) (CustomerInfo, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return CustomerInfo{}, ErrEmptyAddress
	}

	phone = strings.TrimSpace(phone)
	if !phonePattern.MatchString(phone) {
		return CustomerInfo{}, ErrInvalidPhone
	}

	return CustomerInfo{
		address: address,
		phone:   phone,
		note:    strings.TrimSpace(note),
	}, nil
}

// NewCustomerInfoFromRepository - saqlangan CustomerInfo'ni bazadan qayta tiklash uchun,
// validatsiyadan ajratilgan (invariantlar bazadan o'qishda qayta tekshirilmaydi).
func NewCustomerInfoFromRepository(address, phone, note string) CustomerInfo {
	return CustomerInfo{address: address, phone: phone, note: note}
}

func (c CustomerInfo) Address() string { return c.address }
func (c CustomerInfo) Phone() string   { return c.phone }
func (c CustomerInfo) Note() string    { return c.note }
