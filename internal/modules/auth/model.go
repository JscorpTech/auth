package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"first_name"`
	LastName  string `gorm:"last_name"`
	Email     string `gorm:"email;unique"`
	Phone     string `gorm:"phone;unique"`
	Password  string `gorm:"password"`
	Role      string `gorm:"role;default:user"`
}

func (*User) TableName() string {
	return "users"
}
