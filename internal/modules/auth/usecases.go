package auth

import (
	"context"
)

type AuthUsecase interface {
	Login(context.Context) error
	Register(context.Context, *User) (*User, error)
}

type AuthUsecaseImpl struct {
	repo AuthRepository
}

func NewAuthUsecase(repo AuthRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		repo: repo,
	}
}

func (a *AuthUsecaseImpl) Login(cxt context.Context) error {
	return nil
}

func (a *AuthUsecaseImpl) Register(ctx context.Context, user *User) (*User, error) {
	user, err := a.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
