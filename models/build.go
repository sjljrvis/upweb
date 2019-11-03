package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Build struct {
	Base
	CommitHash   string `json:"commit_hash"`
	RepositoryID uint   `json:"repository_id"`
	UserName     string `json:"user_name"`
	email        string `json:"email"`
}
