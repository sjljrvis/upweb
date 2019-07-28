package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

var DB *gorm.DB

// type Resource struct {
// 	gorm.Model

// 	Link        string
// 	Name        string
// 	Author      string
// 	Description string
// 	Tags        pq.StringArray `gorm:"type:varchar(64)[]"`
// }

func Init() {
	db, err := gorm.Open("postgres","host=raja.db.elephantsql.com port=5432 user=mnmvgmco dbname=mnmvgmco sslmode=disable password=BYze89aoMV80jwW6fRwgSKUcDgSiirhB")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	// db.AutoMigrate(&Resource{})
	DB = db
}
