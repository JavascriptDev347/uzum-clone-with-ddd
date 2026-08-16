package media

import "context"

type Uploader interface {
	Upload(ctx context.Context, input UploadInput) (UploadResult, error)
	Delete(ctx context.Context, publicID string) error
}

type UploadInput struct {
	FileName    string // masalan: "category-images/uuid" (kengaytmasiz, Cloudinary o'zi qo'shadi)
	ContentType string // masalan: "image/jpeg"
	Data        []byte
}
type UploadResult struct {
	URL      string
	PublicID string
}
