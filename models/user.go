package models

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	Helper "github.com/sjljrvis/deploynow/helpers"
	fs "github.com/sjljrvis/deploynow/lib/fs"
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

// AfterCreate will create user dir with www-data user.
/*
	//TODO
	1) Create bare repository
	2) Add hook files
	3) Lauch default container
	4) Change owner ship to www-data
*/
func (user *User) AfterCreate(scope *gorm.Scope) (err error) {
	err = fs.CreateDir(path.Join(os.Getenv("ROOT_DIR"), user.UserName))
	if err != nil {
		fmt.Printf("Error Occured")
	}
	return nil
}
