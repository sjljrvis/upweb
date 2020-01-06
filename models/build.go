package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Build struct {
	Base
	CommitHash   string `json:"commit_hash"`
	RepositoryID uint   `json:"repository_id"`
	Email        string `json:"email"`
	UserID       uint   `json:"user_id"`
}
