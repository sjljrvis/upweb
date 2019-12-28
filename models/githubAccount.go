package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type GithubAccount struct {
	Base
	UserID      uint   `gorm:"unique" json:"user_id"`
	AccessToken string `gorm:"unique" json:"access_token"`
	Login       string `gorm:"unique" json:"login"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	AvatarURL   string `json:"avatar_url"`
}
