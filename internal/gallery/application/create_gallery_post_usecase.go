package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

// CreateGalleryPostInput - yangi post yaratish uchun kirish ma'lumotlari. Images'ning
// FileName'lari chaqiruvchi (HTTP handler) tomonidan xavfsiz generatsiya qilingan bo'lishi
// kerak (uuid + kengaytma) - foydalanuvchi yuklagan asl fayl nomi to'g'ridan-to'g'ri
// ishlatilmaydi.
type CreateGalleryPostInput struct {
	Images      []media.UploadInput
	Description string
}

// CreateGalleryPostUseCase - admin uchun: rasm(lar)ni S3'ga yuklaydi va yangi galereya
// postini yaratadi. CreateProductUseCase bilan bir xil yuklash/rollback mantig'i: agar
// keyingi rasmlardan biri yoki post/saqlash muvaffaqiyatsiz bo'lsa, shu chaqiruv ichida
// allaqachon yuklangan rasmlar S3'dan orqaga qaytariladi (hammasi yoki hech narsa).
type CreateGalleryPostUseCase struct {
	repo     domain.GalleryPostRepository
	uploader media.Uploader
}

func NewCreateGalleryPostUseCase(repo domain.GalleryPostRepository, uploader media.Uploader) *CreateGalleryPostUseCase {
	return &CreateGalleryPostUseCase{repo: repo, uploader: uploader}
}

func (uc *CreateGalleryPostUseCase) Execute(ctx context.Context, input CreateGalleryPostInput) (*domain.GalleryPost, error) {
	// Eng arzon tekshiruv birinchi: S3'ga birorta ham rasm yuklashdan oldin rad etamiz.
	if len(input.Images) > domain.MaxGalleryImages {
		return nil, domain.ErrTooManyImages
	}

	uploaded := make([]domain.GalleryImage, 0, len(input.Images))
	rollback := func() {
		for _, img := range uploaded {
			_ = uc.uploader.Delete(ctx, img.PublicID)
		}
	}

	for _, img := range input.Images {
		if err := media.Validate(img, media.DefaultImageRules()); err != nil {
			rollback()
			return nil, err
		}
		result, err := uc.uploader.Upload(ctx, img)
		if err != nil {
			rollback()
			return nil, err
		}
		uploaded = append(uploaded, domain.GalleryImage{URL: result.URL, PublicID: result.PublicID})
	}

	post, err := domain.NewGalleryPost(uploaded, input.Description)
	if err != nil {
		rollback()
		return nil, err
	}

	if err := uc.repo.Save(ctx, post); err != nil {
		rollback()
		return nil, err
	}

	return post, nil
}
