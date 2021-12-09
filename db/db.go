package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (err error) {
	DB, err = gorm.Open(sqlite.Open("cvwo.db"), &gorm.Config{})
	if err != nil {
		return
	}

	return
}
