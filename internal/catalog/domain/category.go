package domain

import (
	"errors"
	"time"
)

var (
	ErrorEmptyCategoryName   = errors.New("catalog: category nomi bo'sh bo'lishi mumkin emas")
	ErrCategoryNotFound      = errors.New("catalog: category topilmadi")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists under the same parent")
)

type Category struct {
	id        string
	name      string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewCategory(id, name string) (*Category, error) {
	if id == "" {
		return nil, ErrEmptyCategoryID
	}
	if name == "" {
		return nil, ErrorEmptyCategoryName
	}

	return &Category{
		id:        id,
		name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		deletedAt: nil,
	}, nil
}

func NewCategoryFromRepository(id, name string, createdAt time.Time, updatedAt time.Time, deletedAt *time.Time) *Category {

	return &Category{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (c *Category) ID() string {
	return c.id
}
func (c *Category) Name() string {
	return c.name
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

func (c *Category) IsDeleted() bool {
	return c.deletedAt != nil
}

func (c *Category) Delete() {
	now := time.Now()
	c.deletedAt = &now
}

func (c *Category) ChangeName(name string) error {
	if name == "" {
		return ErrorEmptyCategoryName
	}
	c.name = name
	c.updatedAt = time.Now()
	return nil
}
