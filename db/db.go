package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sjljrvis/deploynow/log"
	. "github.com/sjljrvis/deploynow/models"
)

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.Info().Msgf("time=%s query=%s", v[2], v[3])
	}
	if v[0] == "log" {
		log.Error().Msgf("Message=%s", v[2])
	}
}

//DB -> sharable DB object to access instance of DB
var DB *gorm.DB

//Init to connect to DB server
func Init() {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=dnow dbname=dnow sslmode=disable password=dnow")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	db.SetLogger(&GormLogger{})
	db.LogMode(true)
	db.AutoMigrate(&Repository{}, &User{}, &Variable{}, &Build{})
	DB = db
}
