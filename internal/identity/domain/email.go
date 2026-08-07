package domain

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidEmail = errors.New("identity: email formati noto'g'ri")

// regexp for validating email format
var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// Email - Value Object.
type Email struct {
	value string
}

// NewEmail - email formatini tekshiradi va kichik harflarga o'giradi (normalize).
func NewEmail(raw string) (Email, error) {
	cleaned := strings.ToLower(strings.TrimSpace(raw))

	if !emailPattern.MatchString(cleaned) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: cleaned}, nil
}

// String - email qiymatini string sifatida qaytaradi.
func (e Email) String() string {
	return e.value
}

// Equals - ikkita Email bir xil qiymatga ega ekanligini tekshiradi (VO tengligi - qiymat bo'yicha).
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
