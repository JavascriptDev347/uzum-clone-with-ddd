package domain

import (
	"errors"
	"testing"
)

func TestNewCustomerInfo(t *testing.T) {
	tests := []struct {
		name    string
		address string
		phone   string
		note    string
		wantErr error
	}{
		{
			name:    "valid, with note",
			address: "Tashkent, Chilonzor 1",
			phone:   "+998901234567",
			note:    "leave at the door",
			wantErr: nil,
		},
		{
			name:    "valid, empty note is allowed",
			address: "Tashkent, Chilonzor 1",
			phone:   "+998901234567",
			note:    "",
			wantErr: nil,
		},
		{
			name:    "empty address",
			address: "",
			phone:   "+998901234567",
			wantErr: ErrEmptyAddress,
		},
		{
			name:    "whitespace-only address",
			address: "   ",
			phone:   "+998901234567",
			wantErr: ErrEmptyAddress,
		},
		{
			name:    "missing country code",
			address: "Tashkent, Chilonzor 1",
			phone:   "901234567",
			wantErr: ErrInvalidPhone,
		},
		{
			name:    "too few digits",
			address: "Tashkent, Chilonzor 1",
			phone:   "+99890123456",
			wantErr: ErrInvalidPhone,
		},
		{
			name:    "too many digits",
			address: "Tashkent, Chilonzor 1",
			phone:   "+9989012345678",
			wantErr: ErrInvalidPhone,
		},
		{
			name:    "non-numeric phone",
			address: "Tashkent, Chilonzor 1",
			phone:   "+998abcdefghi",
			wantErr: ErrInvalidPhone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := NewCustomerInfo(tt.address, tt.phone, tt.note)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}
			if info.Phone() != tt.phone {
				t.Fatalf("expected phone %q, got %q", tt.phone, info.Phone())
			}
		})
	}
}
