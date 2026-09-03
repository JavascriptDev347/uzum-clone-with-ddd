package domain

import (
	"errors"
	"time"
)

var (
	ErrorEmptyCategoryName   = errors.New("catalog: category nomi (uz, eng, ru) bo'sh bo'lishi mumkin emas")
	ErrCategoryNotFound      = errors.New("catalog: category topilmadi")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists under the same parent")
)

type Category struct {
	id      string
	nameUz  string
	nameEng string
	nameRu  string

	imageURL      string
	imagePublicID string
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
}

func NewCategory(id, nameUz, nameEng, nameRu, imageURL, imagePublicID string) (*Category, error) {
	if id == "" {
		return nil, ErrEmptyCategoryID
	}
	if nameUz == "" || nameEng == "" || nameRu == "" {
		return nil, ErrorEmptyCategoryName
	}

	return &Category{
		id:            id,
		nameUz:        nameUz,
		nameEng:       nameEng,
		nameRu:        nameRu,
		imageURL:      imageURL,
		imagePublicID: imagePublicID,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
		deletedAt:     nil,
	}, nil
}

func NewCategoryFromRepository(id, nameUz, nameEng, nameRu, imageURL, imagePublicID string, createdAt time.Time, updatedAt time.Time, deletedAt *time.Time) *Category {

	return &Category{
		id:            id,
		nameUz:        nameUz,
		nameEng:       nameEng,
		nameRu:        nameRu,
		imageURL:      imageURL,
		imagePublicID: imagePublicID,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		deletedAt:     deletedAt,
	}
}

func (c *Category) ID() string {
	return c.id
}
func (c *Category) NameUz() string {
	return c.nameUz
}
func (c *Category) NameEng() string {
	return c.nameEng
}
func (c *Category) NameRu() string {
	return c.nameRu
}
func (c *Category) CreatedAt() time.Time {
	return c.createdAt
}
func (c *Category) UpdatedAt() time.Time {
	return c.updatedAt
}
func (c *Category) DeletedAt() *time.Time {
	return c.deletedAt
}
func (c *Category) ImageURL() string {
	return c.imageURL
}
func (c *Category) ImagePublicID() string {
	return c.imagePublicID
}

func (c *Category) IsDeleted() bool {
	return c.deletedAt != nil
}

func (c *Category) Delete() {
	now := time.Now()
	c.deletedAt = &now
}

func (c *Category) ChangeNames(nameUz, nameEng, nameRu string) error {
	if nameUz == "" || nameEng == "" || nameRu == "" {
		return ErrorEmptyCategoryName
	}
	c.nameUz = nameUz
	c.nameEng = nameEng
	c.nameRu = nameRu
	c.updatedAt = time.Now()
	return nil
}
