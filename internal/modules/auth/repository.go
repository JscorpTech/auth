package auth

import (
	"context"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetID(context.Context, int64) *User
	Create(context.Context, *User) (*User, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (a *AuthRepositoryImpl) GetID(ctx context.Context, id int64) *User {
	return &User{}
}

func (a *AuthRepositoryImpl) Create(ctx context.Context, user *User) (*User, error) {
	if err := a.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
