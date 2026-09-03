package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type UpdateEventUseCase struct {
	repo         domain.EventRepository
	categoryRepo domain.CategoryRepository
}

func NewUpdateEventUseCase(repo domain.EventRepository, categoryRepo domain.CategoryRepository) *UpdateEventUseCase {
	return &UpdateEventUseCase{repo: repo, categoryRepo: categoryRepo}
}

func (uc *UpdateEventUseCase) Execute(ctx context.Context, input UpdateEventInput) error {
	event, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if input.CategoryID != nil {
		if _, err := uc.categoryRepo.FindByID(ctx, *input.CategoryID); err != nil {
			return err
		}
		if err := event.ChangeCategory(*input.CategoryID); err != nil {
			return err
		}
	}

	if input.TitleUz != nil || input.TitleEng != nil || input.TitleRu != nil {
		titleUz, titleEng, titleRu := event.TitleUz(), event.TitleEng(), event.TitleRu()
		if input.TitleUz != nil {
			titleUz = *input.TitleUz
		}
		if input.TitleEng != nil {
			titleEng = *input.TitleEng
		}
		if input.TitleRu != nil {
			titleRu = *input.TitleRu
		}
		if err := event.ChangeTitles(titleUz, titleEng, titleRu); err != nil {
			return err
		}
	}

	if input.EyebrowUz != nil || input.EyebrowEng != nil || input.EyebrowRu != nil {
		eyebrowUz, eyebrowEng, eyebrowRu := event.EyebrowUz(), event.EyebrowEng(), event.EyebrowRu()
		if input.EyebrowUz != nil {
			eyebrowUz = *input.EyebrowUz
		}
		if input.EyebrowEng != nil {
			eyebrowEng = *input.EyebrowEng
		}
		if input.EyebrowRu != nil {
			eyebrowRu = *input.EyebrowRu
		}
		event.ChangeEyebrows(eyebrowUz, eyebrowEng, eyebrowRu)
	}

	if input.SubtitleUz != nil || input.SubtitleEng != nil || input.SubtitleRu != nil {
		subtitleUz, subtitleEng, subtitleRu := event.SubtitleUz(), event.SubtitleEng(), event.SubtitleRu()
		if input.SubtitleUz != nil {
			subtitleUz = *input.SubtitleUz
		}
		if input.SubtitleEng != nil {
			subtitleEng = *input.SubtitleEng
		}
		if input.SubtitleRu != nil {
			subtitleRu = *input.SubtitleRu
		}
		event.ChangeSubtitles(subtitleUz, subtitleEng, subtitleRu)
	}

	if input.CTAUz != nil || input.CTAEng != nil || input.CTARu != nil {
		ctaUz, ctaEng, ctaRu := event.CTAUz(), event.CTAEng(), event.CTARu()
		if input.CTAUz != nil {
			ctaUz = *input.CTAUz
		}
		if input.CTAEng != nil {
			ctaEng = *input.CTAEng
		}
		if input.CTARu != nil {
			ctaRu = *input.CTARu
		}
		event.ChangeCTAs(ctaUz, ctaEng, ctaRu)
	}

	if input.IsRoot != nil {
		event.ChangeIsRoot(*input.IsRoot)
	}

	return uc.repo.Update(ctx, event)
}
