package models

import (
	// "log"
	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)
 
type Repository struct {
	Base
	RepositoryName string                 `gorm:"unique" json:"repositoryName"`
	Language       string                 `json:"language"`
	Path           string                 `json:"path"`
	PathDocker     string                 `json:"pathDocker"`
	Description    string                 `json:"description"`
	State          string                 `json:"state"`
}
 
// BeforeCreate will set a UUID rather than numeric ID.
// func (base *Repo) BeforeCreate(scope *gorm.Scope) error {
// 	log.Printf("I am here")
// 	return nil
//  }