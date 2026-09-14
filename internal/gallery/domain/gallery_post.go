package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// MaxGalleryImages - bitta galereya posti uchun yuklash mumkin bo'lgan eng ko'p rasmlar soni.
const MaxGalleryImages = 3

// GalleryImage - postning bitta rasmi (S3'dagi manzili va obyekt kaliti). catalog'ning
// ProductImage'iga o'xshash: PublicID faqat ichki qatlamlarda (repository, S3 tozalash)
// ishlatiladi va HTTP javobiga chiqarilmaydi (ImageURLs()ga qarang) - lekin
// DeleteGalleryPostUseCase postni o'chirishda S3'dagi obyektlarni ham tozalashi uchun
// saqlanishi shart, faqat URL bilan bu mumkin bo'lmas edi.
type GalleryImage struct {
	URL      string
	PublicID string
}

// GalleryPost - galereyadagi bitta post: rasm(lar) va ixtiyoriy tavsif.
type GalleryPost struct {
	id          uuid.UUID
	images      []GalleryImage
	description string
	createdAt   time.Time
}

// NewGalleryPost - yangi post yaratish uchun. Rasmlar soni shu yerda tekshiriladi - Post o'z
// invariantini o'zi himoya qiladi. Description ixtiyoriy (CustomerInfo.Note bilan bir xil
// qaror: bo'sh bo'lishi mumkin, faqat trim qilinadi, alohida xato talab qilinmaydi).
func NewGalleryPost(images []GalleryImage, description string) (*GalleryPost, error) {
	if len(images) > MaxGalleryImages {
		return nil, ErrTooManyImages
	}

	return &GalleryPost{
		id:          uuid.New(),
		images:      images,
		description: strings.TrimSpace(description),
		createdAt:   time.Now(),
	}, nil
}

// NewGalleryPostFromRepository - saqlangan postni bazadan qayta tiklash uchun, domen
// konstruktoridan ajratilgan (invariantlar bazadan o'qishda qayta tekshirilmaydi).
func NewGalleryPostFromRepository(id uuid.UUID, images []GalleryImage, description string, createdAt time.Time) *GalleryPost {
	return &GalleryPost{
		id:          id,
		images:      images,
		description: description,
		createdAt:   createdAt,
	}
}

func (p *GalleryPost) ID() uuid.UUID { return p.id }

// Images - PublicID bilan birga to'liq rasm ro'yxati (repository va S3 tozalash uchun).
func (p *GalleryPost) Images() []GalleryImage { return p.images }

// ImageURLs - faqat rasm manzillari - HTTP javobida ko'rsatiladigan tashqi shakl.
func (p *GalleryPost) ImageURLs() []string {
	urls := make([]string, 0, len(p.images))
	for _, img := range p.images {
		urls = append(urls, img.URL)
	}
	return urls
}

func (p *GalleryPost) Description() string  { return p.description }
func (p *GalleryPost) CreatedAt() time.Time { return p.createdAt }
