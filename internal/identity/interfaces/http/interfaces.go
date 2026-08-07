package http

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/application"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, input application.RegisterUserInput) (application.RegisterUserOutput, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, input application.LoginUserInput) (application.LoginUserOutput, error)
}

type RefreshUseCase interface {
	Execute(ctx context.Context, input string) (application.LoginUserOutput, error)
}
