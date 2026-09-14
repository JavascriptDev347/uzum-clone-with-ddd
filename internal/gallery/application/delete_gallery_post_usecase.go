package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

// DeleteGalleryPostUseCase - admin uchun: postni bazadan o'chiradi va uning barcha rasmlarini
// S3'dan ham tozalaydi. DeleteProductImageUseCase bilan bir xil tartib: avval baza yozuvi
// o'chiriladi, keyin S3'dan o'chiriladi (xatosi e'tiborsiz qoldiriladi) - shunday qilib baza
// har doim izchil holatda qoladi, va agar S3'dan o'chirish muvaffaqiyatsiz bo'lsa ham
// (masalan tarmoq xatosi), yetim S3 obyekti buzilgan havoladan ko'ra kamroq zarar.
type DeleteGalleryPostUseCase struct {
	repo     domain.GalleryPostRepository
	uploader media.Uploader
}

func NewDeleteGalleryPostUseCase(repo domain.GalleryPostRepository, uploader media.Uploader) *DeleteGalleryPostUseCase {
	return &DeleteGalleryPostUseCase{repo: repo, uploader: uploader}
}

func (uc *DeleteGalleryPostUseCase) Execute(ctx context.Context, id string) error {
	post, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	for _, img := range post.Images() {
		_ = uc.uploader.Delete(ctx, img.PublicID)
	}

	return nil
}
