package application

import (
	"context"
	"errors"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginUserUseCase struct {
	userRepo     domain.UserRepository
	hasher       domain.PasswordHasher
	tokenService domain.TokenService
}

func NewLoginUserUseCase(userRepo domain.UserRepository, hasher domain.PasswordHasher, tokenService domain.TokenService) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {

	// check the email
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return LoginUserOutput{}, err
	}

	// check the user exists
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return LoginUserOutput{}, ErrInvalidCredentials
		}
		return LoginUserOutput{}, ErrInvalidCredentials
	}

	// compare the password
	err = uc.hasher.Compare(user.PasswordHash(), input.Password)
	if err != nil {
		return LoginUserOutput{}, ErrInvalidCredentials
	}

	// generate the access token
	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID(), user.Role())
	if err != nil {
		return LoginUserOutput{}, err
	}

	// generate the refresh token
	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID())
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
