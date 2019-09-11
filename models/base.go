package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

//Base struct to be used for all the dnow models
type Base struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UUID      uuid.UUID `gorm:"type:uuid" json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	return scope.SetColumn("UUID", uuid)
}
