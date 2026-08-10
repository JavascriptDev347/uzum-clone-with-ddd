package http

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
)

type CreateUseCase interface {
	Execute(ctx context.Context, product application.CreateProductInput) (application.CreateProductOutput, error)
}
