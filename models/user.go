package models

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	Helper "github.com/sjljrvis/deploynow/helpers"
	fs "github.com/sjljrvis/deploynow/lib/fs"
	nginx "github.com/sjljrvis/deploynow/lib/nginx"
)

type User struct {
	Base
	UserName      string `gorm:"unique" json:"user_name"`
	Email         string `gorm:"unique" json:"email"`
	Password      string `json:"password"`
	MD5           string `json:"password_md5"`
	Verified      bool   `json:"verified" gorm:"default:false"`
	Repositories  []Repository
	GithubAccount GithubAccount
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeSave(scope *gorm.Scope) error {
	hashed, _ := Helper.HashPassword(user.Password)
	md5, _ := Helper.GetMD5Hash(user.UserName, user.Password)
	scope.SetColumn("MD5", md5[user.UserName])
	return scope.SetColumn("Password", hashed)
}

func (user *User) AfterUpdate(scope *gorm.Scope) error {
	err := nginx.Writehtpasswd(user.UserName, user.MD5)
	err = nginx.Reload()
	if err != nil {
		fmt.Printf("Error Occured")
	}
	return nil
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
	userDir := path.Join(os.Getenv("ROOT_DIR"), user.UserName)
	err = fs.CreateDir(userDir)
	err = nginx.Writehtpasswd(user.UserName, user.MD5)
	err = nginx.Reload()
	if err != nil {
		fmt.Printf("Error Occured")
	}
	return nil
}

func (user *User) BeforeDelete(scope *gorm.Scope) (err error) {
	scope.DB().Where("user_id = ?", user.ID).Delete(Repository{})
	return nil
}

func (user *User) AfterDelete(scope *gorm.Scope) (err error) {
	fmt.Println("Deleting Object")
	return
}
