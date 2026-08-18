package http

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
)

type CreateEventUseCase interface {
	Execute(ctx context.Context, input application.CreateEventInput) (*application.EventOutput, error)
}

type GetEventsUseCase interface {
	Execute(ctx context.Context) ([]*application.EventOutput, error)
}

type GetEventUseCase interface {
	Execute(ctx context.Context, id string) (*application.EventOutput, error)
}

type UpdateEventUseCase interface {
	Execute(ctx context.Context, input application.UpdateEventInput) error
}

type DeleteEventUseCase interface {
	Execute(ctx context.Context, id string) error
}

type GetAllEventsIncludingDeletedUseCase interface {
	Execute(ctx context.Context) ([]*application.EventOutputForAdmin, error)
}
