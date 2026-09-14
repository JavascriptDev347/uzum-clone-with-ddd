// Package money - Money value object'i, avval catalog/domain'da yashagan, endi
// internal/shared/money'ga ko'chirilgan. Uni ordering va (agar kerak bo'lsa) review kabi
// bir nechta bounded context import qilgani uchun, bitta context'ning domenida "tasodifan"
// baham ko'rilishi o'rniga aniq, atayin "shared kernel" sifatida shu yerda joylashadi. Bu
// paket hech qanday bounded context'ga bog'liq emas.
package money

import "errors"

var (
	ErrNegativeAmount   = errors.New("catalog: summa manfiy bo'lishi mumkin emas")
	ErrCurrencyMismatch = errors.New("catalog: valyutalar mos kelmaydi")
	ErrInvalidCurrency  = errors.New("catalog: noto'g'ri valyuta")
)

// amount bu yerda so'mda saqlanadi (tiyin/qism qismlari bilan) - backend tiyinga konvertatsiya qilmaydi,
// qiymat frontenddan qanday kelsa, shundayligicha saqlanadi.
type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}

	if currency == "" {
		return Money{}, ErrInvalidCurrency
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}
