package auth

import (
	"time"

	"gorm.io/gorm"
)

type Role string

var (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleSuper Role = "super"
)

type User struct {
	gorm.Model
	FirstName       string     `gorm:"first_name"`
	LastName        string     `gorm:"last_name"`
	Phone           *string    `gorm:"phone;default:null;uniqueIndex"`
	Email           *string    `gorm:"email;default:null;uniqueIndex"`
	UserName        *string    `gorm:"username;uniqueIndex;default:null"`
	Balance         int        `gorm:"balance"`
	TemplateBalance int        `gorm:"template_balance"`
	Password        string     `gorm:"password"`
	ValidatedAT     *time.Time `gorm:"validated_at"`
	Role            Role       `gorm:"role;default:user"`
}

func (*User) TableName() string {
	return "users"
}

type Otp struct {
	gorm.Model
	Phone string    `gorm:"phone;unique"`
	Code  string    `gorm:"code"`
	Exp   time.Time `gorm:"exp"`
}

func (*Otp) TableName() string {
	return "otp"
}
