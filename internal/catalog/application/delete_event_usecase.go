package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type DeleteEventUseCase struct {
	repo domain.EventRepository
}

func NewDeleteEventUseCase(repo domain.EventRepository) *DeleteEventUseCase {
	return &DeleteEventUseCase{repo: repo}
}

func (uc *DeleteEventUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
