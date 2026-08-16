package media

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	client *cloudinary.Cloudinary
	folder string // masalan: "uzum-clone" - barcha fayllar shu papkada saqlanadi
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
	Folder    string
}

func NewCloudinaryUploader(cfg CloudinaryConfig) (*CloudinaryUploader, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary init error: %w", err)
	}

	return &CloudinaryUploader{
		client: cld,
		folder: cfg.Folder,
	}, nil
}

func (u *CloudinaryUploader) Upload(ctx context.Context, input UploadInput) (UploadResult, error) {
	result, err := u.client.Upload.Upload(ctx, bytes.NewReader(input.Data), uploader.UploadParams{
		Folder:   u.folder,
		PublicID: input.FileName,
	})
	if err != nil {
		return UploadResult{}, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	return UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

func (u *CloudinaryUploader) Delete(ctx context.Context, publicID string) error {
	_, err := u.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}
func (u *CloudinaryUploader) extractPublicID(url string) string {
	// "/upload/" dan keyingi qismni olamiz
	uploadIdx := strings.Index(url, "/upload/")
	if uploadIdx == -1 {
		return ""
	}

	afterUpload := url[uploadIdx+len("/upload/"):]

	// versiya prefiksini olib tashlaymiz (masalan "v1234567890/")
	parts := strings.SplitN(afterUpload, "/", 2)
	var pathWithExt string
	if len(parts) == 2 && strings.HasPrefix(parts[0], "v") && isNumeric(parts[0][1:]) {
		pathWithExt = parts[1]
	} else {
		pathWithExt = afterUpload
	}

	// fayl kengaytmasini (.jpg, .png va h.k.) olib tashlaymiz
	if dotIdx := strings.LastIndex(pathWithExt, "."); dotIdx != -1 {
		pathWithExt = pathWithExt[:dotIdx]
	}

	return pathWithExt
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
