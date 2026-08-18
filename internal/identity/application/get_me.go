package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
)

type GetMeUseCase struct {
	userRepo domain.UserRepository
}

func NewGetMeUseCase(userRepo domain.UserRepository) *GetMeUseCase {
	return &GetMeUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetMeUseCase) Execute(ctx context.Context, userID string) (*GetMeOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetMeOutput{
		UserID: user.ID(),
		Email:  user.Email().String(),
		Role:   string(user.Role()),
	}, nil
}
