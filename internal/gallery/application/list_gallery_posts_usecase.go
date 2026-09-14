package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
)

// ListGalleryPostsUseCase - galereya postlarini sahifalab olish (ochiq, autentifikatsiya
// shart emas).
type ListGalleryPostsUseCase struct {
	repo domain.GalleryPostRepository
}

func NewListGalleryPostsUseCase(repo domain.GalleryPostRepository) *ListGalleryPostsUseCase {
	return &ListGalleryPostsUseCase{repo: repo}
}

func (uc *ListGalleryPostsUseCase) Execute(ctx context.Context, page, pageSize int) ([]*domain.GalleryPost, int64, error) {
	page, pageSize = NormalizeGalleryPagination(page, pageSize)
	return uc.repo.FindAll(ctx, page, pageSize)
}
