package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type GetEventUseCase struct {
	repo domain.EventRepository
}

func NewGetEventUseCase(repo domain.EventRepository) *GetEventUseCase {
	return &GetEventUseCase{repo: repo}
}

func (uc *GetEventUseCase) Execute(ctx context.Context, id string) (*EventOutput, error) {
	event, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToEventOutput(event), nil
}
