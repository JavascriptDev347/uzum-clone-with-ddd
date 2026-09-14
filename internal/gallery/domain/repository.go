package domain

import "context"

type GalleryPostRepository interface {
	Save(ctx context.Context, post *GalleryPost) error

	FindAll(ctx context.Context, page, pageSize int) ([]*GalleryPost, int64, error)

	// Delete - postni bazadan o'chiradi va o'chirilgan postni qaytaradi - chaqiruvchi
	// (DeleteGalleryPostUseCase) shu orqali rasmlarning PublicID'larini bilib, ularni S3'dan
	// ham tozalay oladi. Topilmasa ErrGalleryPostNotFound qaytaradi.
	Delete(ctx context.Context, id string) (*GalleryPost, error)
}
