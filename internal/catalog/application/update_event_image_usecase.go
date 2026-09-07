package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

type UpdateEventImageUseCase struct {
	repo     domain.EventRepository
	uploader media.Uploader
}

func NewUpdateEventImageUseCase(repo domain.EventRepository, uploader media.Uploader) *UpdateEventImageUseCase {
	return &UpdateEventImageUseCase{repo: repo, uploader: uploader}
}

// Execute - eventning rasmini yangisiga almashtiradi va eski rasmni uploader'dan o'chiradi.
func (uc *UpdateEventImageUseCase) Execute(ctx context.Context, input UpdateEventImageInput) (*EventOutput, error) {
	event, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := media.Validate(input.Image, media.DefaultImageRules()); err != nil {
		return nil, err
	}

	uploadResult, err := uc.uploader.Upload(ctx, input.Image)
	if err != nil {
		return nil, err
	}

	oldPublicID := event.ImagePublicID()
	event.ChangeImage(uploadResult.URL, uploadResult.PublicID)

	if err := uc.repo.Update(ctx, event); err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	if oldPublicID != "" {
		_ = uc.uploader.Delete(ctx, oldPublicID)
	}

	return ToEventOutput(event, LangUZ), nil
}
