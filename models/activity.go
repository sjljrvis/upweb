package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Activity struct {
	Base
	Type         string `json:"type"`
	Email        string `json:"email"`
	UserID       uint   `json:"user_id"`
	RepositoryId uint   `json:"repository_id"`
	// MetaData     interface{} `json:"meta_data"`
}

func (activity *Activity) Log(db *gorm.DB, repository Repository, user User) {
	activity.Type = "Initial release"
	activity.Email = user.Email
	activity.UserID = user.ID
	activity.RepositoryId = repository.ID
	db.Save(&activity)
}
