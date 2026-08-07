package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
)

type RefreshTokenUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

func NewRefreshTokenUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (LoginUserOutput, error) {
	// separate the user id from the token
	id, err := uc.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return LoginUserOutput{}, err
	}

	// find the user that available yet
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return LoginUserOutput{}, err
	}

	// generate new tokens
	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID(), user.Role())
	if err != nil {
		return LoginUserOutput{}, err
	}

	refreshToken, err = uc.tokenService.GenerateRefreshToken(user.ID())
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
