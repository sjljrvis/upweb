package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	Helper "github.com/sjljrvis/deploynow/helpers"
)

type User struct {
	Base
	UserName     string `gorm:"unique" json:"user_name"`
	Email        string `gorm:"unique" json:"email"`
	Password     string `json:"password"`
	Repositories Repository
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeSave(scope *gorm.Scope) error {
	hashed, _ := Helper.HashPassword(user.Password)
	return scope.SetColumn("Password", hashed)
}
