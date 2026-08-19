package domain

import "errors"

var (
	ErrInvalidPackagingType     = errors.New("catalog: qadoqlash turi noto'g'ri (bucket, box yoki vase bo'lishi kerak)")
	ErrInvalidFreshnessLifespan = errors.New("catalog: saqlanish muddati 1 dan 7 kungacha bo'lishi kerak")
)

// PackagingType - mahsulot qanday qadoqlanganini bildiradi.
type PackagingType string

const (
	PackagingBucket PackagingType = "bucket"
	PackagingBox    PackagingType = "box"
	PackagingVase   PackagingType = "vase"
)

func (p PackagingType) IsValid() bool {
	switch p {
	case PackagingBucket, PackagingBox, PackagingVase:
		return true
	default:
		return false
	}
}

func (p PackagingType) String() string {
	return string(p)
}

// FreshnessLifespan - buket necha kun saqlanishi (kunlarda), 1 dan 7 gacha bo'ladi.
type FreshnessLifespan int

const DefaultFreshnessLifespan FreshnessLifespan = 1

func (f FreshnessLifespan) IsValid() bool {
	return f >= 1 && f <= 7
}
