package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/google/uuid"
)

type CreateEventUseCase struct {
	repo         domain.EventRepository
	categoryRepo domain.CategoryRepository
	uploader     media.Uploader
}

func NewCreateEventUseCase(repo domain.EventRepository, categoryRepo domain.CategoryRepository, uploader media.Uploader) *CreateEventUseCase {
	return &CreateEventUseCase{repo: repo, categoryRepo: categoryRepo, uploader: uploader}
}

func (uc *CreateEventUseCase) Execute(ctx context.Context, input CreateEventInput) (*EventOutput, error) {
	if _, err := uc.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	if err := media.Validate(input.Image, media.DefaultImageRules()); err != nil {
		return nil, err
	}

	uploadResult, err := uc.uploader.Upload(ctx, input.Image)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	event, err := domain.NewEvent(domain.NewEventParams{
		ID:            id,
		EyebrowUz:     input.EyebrowUz,
		EyebrowEng:    input.EyebrowEng,
		EyebrowRu:     input.EyebrowRu,
		TitleUz:       input.TitleUz,
		TitleEng:      input.TitleEng,
		TitleRu:       input.TitleRu,
		SubtitleUz:    input.SubtitleUz,
		SubtitleEng:   input.SubtitleEng,
		SubtitleRu:    input.SubtitleRu,
		CTAUz:         input.CTAUz,
		CTAEng:        input.CTAEng,
		CTARu:         input.CTARu,
		ImageURL:      uploadResult.URL,
		ImagePublicID: uploadResult.PublicID,
		CategoryID:    input.CategoryID,
		IsRoot:        input.IsRoot,
	})
	if err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	if err := uc.repo.Save(ctx, event); err != nil {
		_ = uc.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	return ToEventOutput(event, LangUZ), nil
}
