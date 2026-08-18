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

	if input.Eyebrow != nil || input.Title != nil || input.Subtitle != nil || input.CTA != nil {
		if err := event.ChangeContent(input.Eyebrow, input.Title, input.Subtitle, input.CTA); err != nil {
			return err
		}
	}

	if input.IsRoot != nil {
		event.ChangeIsRoot(*input.IsRoot)
	}

	return uc.repo.Update(ctx, event)
}
