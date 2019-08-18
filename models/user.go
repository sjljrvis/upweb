package models

import (
	"github.com/jinzhu/gorm"
	Helper "github.com/sjljrvis/deploynow/helpers"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)
 
type User struct {
	Base
	UserName string                 `gorm:"unique" json:"user_name"`
	Email    string                 `gorm:"unique" json:"email"`
	Password string                 `json:"password"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeSave(scope *gorm.Scope) error {
	hashed, _ := Helper.HashPassword(user.Password)
	return scope.SetColumn("Password", hashed)
 }