package auth

import (
	"context"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetID(context.Context, int64) *User
	Create(context.Context, *User) (*User, error)
	IsExists(context.Context, string) bool
	GetByPhone(context.Context, string) (*User, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (a *AuthRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	if err := a.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepositoryImpl) GetID(ctx context.Context, id int64) *User {
	return &User{}
}

func (a *AuthRepositoryImpl) IsExists(ctx context.Context, phone string) bool {
	var count int64
	a.db.WithContext(ctx).Model(&User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

func (a *AuthRepositoryImpl) Create(ctx context.Context, user *User) (*User, error) {
	if err := a.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
