package application

import (
	"context"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/google/uuid"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type RegisterUserUseCase struct {
	userRepo domain.UserRepository
	hasher   domain.PasswordHasher
}

func NewRegisterUserUseCase(userRepo domain.UserRepository, hasher domain.PasswordHasher) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {

	// chec email
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	// if user exists, return error that user already exists
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return RegisterUserOutput{}, err
		}
	}
	if user != nil {
		return RegisterUserOutput{}, ErrEmailAlreadyExists
	}

	// hash the password
	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	newUser, err := domain.NewUser(uuid.New().String(), email, hashedPassword, domain.RoleCustomer, time.Now())
	if err != nil {
		return RegisterUserOutput{}, err
	}

	// save the user
	err = uc.userRepo.Save(ctx, newUser)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		UserID: newUser.ID(),
		Email:  newUser.Email().String(),
	}, nil
}
