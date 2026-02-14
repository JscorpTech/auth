package auth

import (
	"context"
	"errors"
	"time"

	"github.com/JscorpTech/auth/internal/config"
	"github.com/JscorpTech/auth/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase interface {
	Login(context.Context, string, string) (*User, error)
	Register(context.Context, *User) (*User, error)
	IsExists(context.Context, string) bool
	AccessToken(*User) string
	RefreshToken(*User) string
}

type AuthUsecaseImpl struct {
	repo AuthRepository
	cfg  *config.Config
}

func NewAuthUsecase(repo AuthRepository, cfg *config.Config) AuthUsecase {
	return &AuthUsecaseImpl{
		repo: repo,
		cfg:  cfg,
	}
}

func (a *AuthUsecaseImpl) Login(ctx context.Context, phone string, password string) (*User, error) {
	user, err := a.repo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if res := utils.CheckPasswordHash(password, user.Password); !res {
		return nil, errors.New("Invalid password")
	}
	return user, nil
}

func (a *AuthUsecaseImpl) IsExists(ctx context.Context, phone string) bool {
	return a.repo.IsExists(ctx, phone)
}

func (a *AuthUsecaseImpl) Register(ctx context.Context, user *User) (*User, error) {
	if a.repo.IsExists(ctx, user.Phone) {
		return nil, ErrUserAlreadyExists
	}
	user, err := a.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthUsecaseImpl) AccessToken(user *User) string {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"type":    "access",
	}
	token, err := utils.CreateJWT(claims, a.cfg.PrivateKey)
	if err != nil {
		return ""
	}
	return token
}

func (a *AuthUsecaseImpl) RefreshToken(user *User) string {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"type":    "refresh",
	}
	token, err := utils.CreateJWT(claims, a.cfg.PrivateKey)
	if err != nil {
		return ""
	}
	return token
}
