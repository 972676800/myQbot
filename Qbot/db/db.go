package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("tc.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

}
