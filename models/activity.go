package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Activity struct {
	Base
	Info         string `json:"info" gorm:"default:'initial release'"`
	Type         string `json:"type" gorm:"default:'build'"`
	Email        string `json:"email"`
	UserID       uint   `json:"user_id"`
	RepositoryId uint   `json:"repository_id"`
	// MetaData     interface{} `json:"meta_data"`
}

func (activity *Activity) Log(db *gorm.DB, repository_id uint, user User, activity_type , info string) {
	activity.Type = activity_type
	activity.Info = info
	activity.Email = user.Email
	activity.UserID = user.ID
	activity.RepositoryId = repository_id
	db.Save(&activity)
}
