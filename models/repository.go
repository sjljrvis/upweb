package models

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/phayes/freeport"
	uuid "github.com/satori/go.uuid"
	container "github.com/sjljrvis/deploynow/lib/container"
	digitalocean "github.com/sjljrvis/deploynow/lib/digitalocean"
	fs "github.com/sjljrvis/deploynow/lib/fs"
	git "github.com/sjljrvis/deploynow/lib/git"
	nginx "github.com/sjljrvis/deploynow/lib/nginx"
)

type Repository struct {
	Base
	RepositoryName string     `gorm:"unique" json:"repository_name"`
	Language       string     `json:"language"`
	Path           string     `json:"path"`
	PathDocker     string     `json:"path_docker"`
	Description    string     `json:"description"`
	State          string     `json:"state" gorm:"default:'stopped'"`
	UserID         uint       `json:"user_id"`
	UserName       string     `json:"user_name"`
	ContainerID    string     `json:"container_id" gorm:"default:'0'"`
	DNSID          string     `json:"dns_id"`
	Variables      []Variable `json:"variables"`
	GithubLinked   bool       `json:"github_linked"`
	GithubURL      string     `json:"github_url"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (repo *Repository) BeforeCreate(scope *gorm.Scope) (err error) {
	repositoryPath := path.Join(os.Getenv("ROOT_DIR"), repo.UserName, repo.RepositoryName)
	repositoryPathDocker := path.Join(os.Getenv("ROOT_DIR"), repo.UserName, repo.RepositoryName+"_docker")
	uuid, _ := uuid.NewV4()
	scope.SetColumn("UUID", uuid)
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
	container_id := container.GenerateDefault(repo.RepositoryName, port)
	dns_id, err := digitalocean.CreateDNS(repo.RepositoryName)
	if err != nil {
		fmt.Printf("Error Occured")
	}
	scope.DB().Model(repo).Update("State", "running")
	scope.DB().Model(repo).Update("ContainerID", container_id)
	scope.DB().Model(repo).Update("DNSID", dns_id)
	return
}

func (repository *Repository) BeforeDelete(scope *gorm.Scope) (err error) {
	scope.DB().Where("repository_id = ?", repository.ID).Delete(Variable{})
	scope.DB().Where("repository_id = ?", repository.ID).Delete(Build{})
	scope.DB().Where("repository_id = ?", repository.ID).Delete(Activity{})
	return nil
}

// AfterCreate will set a UUID rather than numeric ID.
/*
	//TODO
	1) Clean associated dirs
	2) Add hook files
	3) Lauch default container
	4) Change owner ship to www-data
*/
func (repository *Repository) AfterDelete(scope *gorm.Scope) (err error) {
	err = fs.RemoveDir(repository.Path)
	err = fs.RemoveDir(repository.PathDocker)
	container.Stop(repository.ContainerID)
	container.Remove(repository.ContainerID)
	digitalocean.RemoveDNS(repository.DNSID)
	return nil
}
