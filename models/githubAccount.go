package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type GithubAccount struct {
	Base
	AccessToken string `gorm:"unique" json:"access_token"`
	GID         string `gorm:"unique" json:"g_id"`
	Login       string `json:"login"`
	NodeId      string `json:"node_id"`
	Url         string `json:"url"`
	AvatarUrl   string `json:"avatar_url"`
}
