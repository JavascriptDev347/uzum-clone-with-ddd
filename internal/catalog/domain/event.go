package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyEventID    = errors.New("catalog: event ID bo'sh bo'lishi mumkin emas")
	ErrEmptyEventTitle = errors.New("catalog: event title (uz, eng, ru) bo'sh bo'lishi mumkin emas")
	ErrEventNotFound   = errors.New("catalog: event topilmadi")
)

type Event struct {
	id string

	eyebrowUz  string
	eyebrowEng string
	eyebrowRu  string

	titleUz  string
	titleEng string
	titleRu  string

	subtitleUz  string
	subtitleEng string
	subtitleRu  string

	ctaUz  string
	ctaEng string
	ctaRu  string

	imageURL      string
	imagePublicID string
	categoryID    string
	isRoot        bool
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
}

// NewEventParams - yangi event yaratish uchun kerakli ma'lumotlar.
type NewEventParams struct {
	ID string

	EyebrowUz  string
	EyebrowEng string
	EyebrowRu  string

	TitleUz  string
	TitleEng string
	TitleRu  string

	SubtitleUz  string
	SubtitleEng string
	SubtitleRu  string

	CTAUz  string
	CTAEng string
	CTARu  string

	ImageURL      string
	ImagePublicID string
	CategoryID    string
	IsRoot        bool
}

func NewEvent(p NewEventParams) (*Event, error) {
	if p.ID == "" {
		return nil, ErrEmptyEventID
	}
	if p.TitleUz == "" || p.TitleEng == "" || p.TitleRu == "" {
		return nil, ErrEmptyEventTitle
	}
	if p.CategoryID == "" {
		return nil, ErrEmptyCategoryID
	}

	now := time.Now()
	return &Event{
		id:            p.ID,
		eyebrowUz:     p.EyebrowUz,
		eyebrowEng:    p.EyebrowEng,
		eyebrowRu:     p.EyebrowRu,
		titleUz:       p.TitleUz,
		titleEng:      p.TitleEng,
		titleRu:       p.TitleRu,
		subtitleUz:    p.SubtitleUz,
		subtitleEng:   p.SubtitleEng,
		subtitleRu:    p.SubtitleRu,
		ctaUz:         p.CTAUz,
		ctaEng:        p.CTAEng,
		ctaRu:         p.CTARu,
		imageURL:      p.ImageURL,
		imagePublicID: p.ImagePublicID,
		categoryID:    p.CategoryID,
		isRoot:        p.IsRoot,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// EventFromRepositoryParams - saqlangan eventni bazadan qayta tiklash uchun.
type EventFromRepositoryParams struct {
	ID string

	EyebrowUz  string
	EyebrowEng string
	EyebrowRu  string

	TitleUz  string
	TitleEng string
	TitleRu  string

	SubtitleUz  string
	SubtitleEng string
	SubtitleRu  string

	CTAUz  string
	CTAEng string
	CTARu  string

	ImageURL      string
	ImagePublicID string
	CategoryID    string
	IsRoot        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func NewEventFromRepository(p EventFromRepositoryParams) *Event {
	return &Event{
		id:            p.ID,
		eyebrowUz:     p.EyebrowUz,
		eyebrowEng:    p.EyebrowEng,
		eyebrowRu:     p.EyebrowRu,
		titleUz:       p.TitleUz,
		titleEng:      p.TitleEng,
		titleRu:       p.TitleRu,
		subtitleUz:    p.SubtitleUz,
		subtitleEng:   p.SubtitleEng,
		subtitleRu:    p.SubtitleRu,
		ctaUz:         p.CTAUz,
		ctaEng:        p.CTAEng,
		ctaRu:         p.CTARu,
		imageURL:      p.ImageURL,
		imagePublicID: p.ImagePublicID,
		categoryID:    p.CategoryID,
		isRoot:        p.IsRoot,
		createdAt:     p.CreatedAt,
		updatedAt:     p.UpdatedAt,
		deletedAt:     p.DeletedAt,
	}
}

func (e *Event) ID() string            { return e.id }
func (e *Event) EyebrowUz() string     { return e.eyebrowUz }
func (e *Event) EyebrowEng() string    { return e.eyebrowEng }
func (e *Event) EyebrowRu() string     { return e.eyebrowRu }
func (e *Event) TitleUz() string       { return e.titleUz }
func (e *Event) TitleEng() string      { return e.titleEng }
func (e *Event) TitleRu() string       { return e.titleRu }
func (e *Event) SubtitleUz() string    { return e.subtitleUz }
func (e *Event) SubtitleEng() string   { return e.subtitleEng }
func (e *Event) SubtitleRu() string    { return e.subtitleRu }
func (e *Event) CTAUz() string         { return e.ctaUz }
func (e *Event) CTAEng() string        { return e.ctaEng }
func (e *Event) CTARu() string         { return e.ctaRu }
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

func (e *Event) ChangeTitles(titleUz, titleEng, titleRu string) error {
	if titleUz == "" || titleEng == "" || titleRu == "" {
		return ErrEmptyEventTitle
	}
	e.titleUz = titleUz
	e.titleEng = titleEng
	e.titleRu = titleRu
	e.updatedAt = time.Now()
	return nil
}

func (e *Event) ChangeEyebrows(eyebrowUz, eyebrowEng, eyebrowRu string) {
	e.eyebrowUz = eyebrowUz
	e.eyebrowEng = eyebrowEng
	e.eyebrowRu = eyebrowRu
	e.updatedAt = time.Now()
}

func (e *Event) ChangeSubtitles(subtitleUz, subtitleEng, subtitleRu string) {
	e.subtitleUz = subtitleUz
	e.subtitleEng = subtitleEng
	e.subtitleRu = subtitleRu
	e.updatedAt = time.Now()
}

func (e *Event) ChangeCTAs(ctaUz, ctaEng, ctaRu string) {
	e.ctaUz = ctaUz
	e.ctaEng = ctaEng
	e.ctaRu = ctaRu
	e.updatedAt = time.Now()
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
