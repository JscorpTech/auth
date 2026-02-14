package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"first_name"`
	LastName  string `gorm:"last_name"`
	Email     string `gorm:"email"`
	Phone     string `gorm:"phone"`
	Password  string `gorm:"password"`
}

func (*User) TableName() string {
	return "users"
}
