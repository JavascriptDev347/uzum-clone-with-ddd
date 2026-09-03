package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetAllEventsIncludingDeletedUseCase struct {
	repo domain.EventRepository
}

func NewGetAllEventsIncludingDeletedUseCase(repo domain.EventRepository) *GetAllEventsIncludingDeletedUseCase {
	return &GetAllEventsIncludingDeletedUseCase{repo: repo}
}

func (uc *GetAllEventsIncludingDeletedUseCase) Execute(ctx context.Context) ([]*EventOutputForAdmin, error) {
	events, err := uc.repo.FindAllIncludingDeleted(ctx)
	if err != nil {
		return nil, err
	}

	return ToEventOutputsForAdmin(events), nil
}
