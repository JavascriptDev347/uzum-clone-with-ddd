package domain

import "errors"

var (
	ErrNegativeAmount   = errors.New("catalog: summa manfiy bo'lishi mumkin emas")
	ErrCurrencyMismatch = errors.New("catalog: valyutalar mos kelmaydi")
	ErrInvalidCurrency  = errors.New("catalog: noto'g'ri valyuta")
)

// amount bu yerda tiyinlarda saqlanadi
type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
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

func (m *Money) Amount() int64 {
	return m.amount
}

func (m *Money) Currency() string {
	return m.currency
}
