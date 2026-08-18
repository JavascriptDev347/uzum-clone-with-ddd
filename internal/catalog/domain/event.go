package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyEventID    = errors.New("catalog: event ID bo'sh bo'lishi mumkin emas")
	ErrEmptyEventTitle = errors.New("catalog: event title bo'sh bo'lishi mumkin emas")
	ErrEventNotFound   = errors.New("catalog: event topilmadi")
)

type Event struct {
	id            string
	eyebrow       string
	title         string
	subtitle      string
	cta           string
	imageURL      string
	imagePublicID string
	categoryID    string
	isRoot        bool
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
}

func NewEvent(id, eyebrow, title, subtitle, cta, imageURL, imagePublicID, categoryID string, isRoot bool) (*Event, error) {
	if id == "" {
		return nil, ErrEmptyEventID
	}
	if title == "" {
		return nil, ErrEmptyEventTitle
	}
	if categoryID == "" {
		return nil, ErrEmptyCategoryID
	}

	now := time.Now()
	return &Event{
		id:            id,
		eyebrow:       eyebrow,
		title:         title,
		subtitle:      subtitle,
		cta:           cta,
		imageURL:      imageURL,
		imagePublicID: imagePublicID,
		categoryID:    categoryID,
		isRoot:        isRoot,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func NewEventFromRepository(id, eyebrow, title, subtitle, cta, imageURL, imagePublicID, categoryID string, isRoot bool, createdAt, updatedAt time.Time, deletedAt *time.Time) *Event {
	return &Event{
		id:            id,
		eyebrow:       eyebrow,
		title:         title,
		subtitle:      subtitle,
		cta:           cta,
		imageURL:      imageURL,
		imagePublicID: imagePublicID,
		categoryID:    categoryID,
		isRoot:        isRoot,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		deletedAt:     deletedAt,
	}
}

func (e *Event) ID() string            { return e.id }
func (e *Event) Eyebrow() string       { return e.eyebrow }
func (e *Event) Title() string         { return e.title }
func (e *Event) Subtitle() string      { return e.subtitle }
func (e *Event) CTA() string           { return e.cta }
func (e *Event) ImageURL() string      { return e.imageURL }
func (e *Event) ImagePublicID() string { return e.imagePublicID }
func (e *Event) CategoryID() string    { return e.categoryID }
func (e *Event) IsRoot() bool          { return e.isRoot }
func (e *Event) CreatedAt() time.Time  { return e.createdAt }
func (e *Event) UpdatedAt() time.Time  { return e.updatedAt }
func (e *Event) DeletedAt() *time.Time { return e.deletedAt }

func (e *Event) IsDeleted() bool {
	return e.deletedAt != nil
}

func (e *Event) ChangeContent(eyebrow, title, subtitle, cta *string) error {
	if title != nil {
		if *title == "" {
			return ErrEmptyEventTitle
		}
		e.title = *title
	}
	if eyebrow != nil {
		e.eyebrow = *eyebrow
	}
	if subtitle != nil {
		e.subtitle = *subtitle
	}
	if cta != nil {
		e.cta = *cta
	}
	e.updatedAt = time.Now()
	return nil
}

func (e *Event) ChangeCategory(categoryID string) error {
	if categoryID == "" {
		return ErrEmptyCategoryID
	}
	e.categoryID = categoryID
	e.updatedAt = time.Now()
	return nil
}

func (e *Event) ChangeIsRoot(isRoot bool) {
	e.isRoot = isRoot
	e.updatedAt = time.Now()
}

func (e *Event) Delete() {
	now := time.Now()
	e.deletedAt = &now
}
