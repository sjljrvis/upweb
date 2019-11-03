package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/sjljrvis/deploynow/models"
)

//DB -> sharable DB object to access instance of DB
var DB *gorm.DB

//Init to connect to DB server
func Init() {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=dnow dbname=dnow sslmode=disable password=dnow")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	db.LogMode(true)
	db.AutoMigrate(&Repository{}, &User{}, &Variable{})
	DB = db
}
