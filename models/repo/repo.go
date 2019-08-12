package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)
 
type Repo struct {
	gorm.Model
	ID       			 uuid.UUID 							`gorm:"type:uuid;primary_key;"`
	RepositoryName string                 `gorm:"unique" json:"repositoryName"`
	Language       string                 `json:"language"`
	Path           string                 `json:"path"`
	PathDocker     string                 `json:"pathDocker"`
	Date           time.Time              `json:"date"`
	Description    string                 `json:"description"`
	State          string                 `json:"state"`
	Github         map[string]interface{} `json:"github"`
}
 
// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}
	return scope.SetColumn("ID", uuid)
 }
 
// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
// func DBMigrate(db *gorm.DB) *gorm.DB {
// 	db.AutoMigrate(&Employee{})
// 	return db
// }

