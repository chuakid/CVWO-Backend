package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (err error) {
	var DATABASE_URL string = os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		return
	}

	return
}
