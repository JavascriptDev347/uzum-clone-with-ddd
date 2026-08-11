package domain

import (
	"errors"
	"time"
)

var (
	ErrorEmptyCategoryName = errors.New("catalog: category nomi bo'sh bo'lishi mumkin emas")
	ErrCategoryNotFound    = errors.New("catalog: category topilmadi")
)

type Category struct {
	id        string
	name      string
	parentID  *string
	createdAt time.Time
	deletedAt *time.Time
}

func NewCategory(id, name string, parentID *string) (*Category, error) {
	if id == "" {
		return nil, ErrEmptyCategoryID
	}
	if name == "" {
		return nil, ErrorEmptyCategoryName
	}

	return &Category{
		id:        id,
		name:      name,
		parentID:  parentID,
		createdAt: time.Now(),
		deletedAt: nil,
	}, nil
}

func NewCategoryFromRepository(id, name string, parentID *string, createdAt time.Time, deletedAt *time.Time) *Category {

	return &Category{
		id:        id,
		name:      name,
		parentID:  parentID,
		createdAt: createdAt,
		deletedAt: deletedAt,
	}
}

func (c *Category) ID() string {
	return c.id
}
func (c *Category) Name() string {
	return c.name
}
func (c *Category) ParentID() *string {
	return c.parentID
}
func (c *Category) CreatedAt() time.Time {
	return c.createdAt
}
func (c *Category) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *Category) IsRoot() bool {
	return c.parentID == nil
}

func (c *Category) IsDeleted() bool {
	return c.deletedAt != nil
}

func (c *Category) Delete() {
	now := time.Now()
	c.deletedAt = &now
}
