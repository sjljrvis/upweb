package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Variable struct {
	Base
	Key          string `json:"key"`
	Value        string `json:"value"`
	RepositoryID uint   `json:"repository_id"`
}
