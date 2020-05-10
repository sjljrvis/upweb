package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type EmailToken struct {
	Base
	UserID uint `json:"user_id"`
	status bool `json:"status"`
}
