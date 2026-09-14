package http

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
)

// GalleryPostResponse - postning HTTP javob shakli. domain.GalleryPost to'g'ridan-to'g'ri
// serializatsiya qilinmaydi (private field'lari bor) - ichki PublicID'lar ham tashqariga
// chiqarilmaydi, faqat rasm manzillari (ImageURLs) beriladi.
type GalleryPostResponse struct {
	ID          string    `json:"id"`
	ImageURLs   []string  `json:"image_urls"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToGalleryPostResponse(post *domain.GalleryPost) GalleryPostResponse {
	return GalleryPostResponse{
		ID:          post.ID().String(),
		ImageURLs:   post.ImageURLs(),
		Description: post.Description(),
		CreatedAt:   post.CreatedAt(),
	}
}

func ToGalleryPostResponses(posts []*domain.GalleryPost) []GalleryPostResponse {
	responses := make([]GalleryPostResponse, 0, len(posts))
	for _, p := range posts {
		responses = append(responses, ToGalleryPostResponse(p))
	}
	return responses
}
