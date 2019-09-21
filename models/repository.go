package models

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/phayes/freeport"
	container "github.com/sjljrvis/deploynow/lib/container"
	fs "github.com/sjljrvis/deploynow/lib/fs"
	git "github.com/sjljrvis/deploynow/lib/git"
	nginx "github.com/sjljrvis/deploynow/lib/nginx"
)

type Repository struct {
	Base
	RepositoryName string `gorm:"unique" json:"repository_name"`
	Language       string `json:"language"`
	Path           string `json:"path"`
	PathDocker     string `json:"path_docker"`
	Description    string `json:"description"`
	State          string `json:"state" default:"stopped"`
	UserID         uint   `json:"user_id"`
	UserName       string `json::"user_name"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (repo *Repository) BeforeCreate(scope *gorm.Scope) (err error) {
	repositoryPath := path.Join(os.Getenv("ROOT_DIR"), repo.UserName, repo.RepositoryName)
	repositoryPathDocker := path.Join(os.Getenv("ROOT_DIR"), repo.UserName, repo.RepositoryName+"_docker")
	scope.SetColumn("Path", repositoryPath)
	scope.SetColumn("PathDocker", repositoryPathDocker)
	return
}

// AfterCreate will set a UUID rather than numeric ID.
/*
	//TODO
	1) Create bare repository
	2) Add hook files
	3) Lauch default container
	4) Change owner ship to www-data
*/
func (repo *Repository) AfterCreate(scope *gorm.Scope) (err error) {
	port, err := freeport.GetFreePort()
	err = fs.CreateDir(repo.Path)
	err = fs.CreateDir(repo.PathDocker)
	err = git.InitBare(repo.Path)
	err = git.CreateHooks(repo.Path)
	nginx.WriteConfig(repo.RepositoryName, strconv.Itoa(port))
	nginx.Symlink(repo.RepositoryName)
	container.GenerateDefault(repo.RepositoryName, port)
	if err != nil {
		fmt.Printf("Error Occured")
	}
	return scope.SetColumn("State", "stopped")
}
