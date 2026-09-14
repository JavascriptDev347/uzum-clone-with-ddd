package domain

import "errors"

var (
	ErrTooManyImages       = errors.New("gallery: bitta post uchun eng ko'pi bilan 3 ta rasm yuklash mumkin")
	ErrGalleryPostNotFound = errors.New("gallery: post topilmadi")
)
