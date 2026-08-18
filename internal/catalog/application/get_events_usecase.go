package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetEventsUseCase struct {
	repo domain.EventRepository
}

func NewGetEventsUseCase(repo domain.EventRepository) *GetEventsUseCase {
	return &GetEventsUseCase{repo: repo}
}

func (uc *GetEventsUseCase) Execute(ctx context.Context) ([]*EventOutput, error) {
	events, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return ToEventOutputs(events), nil
}
