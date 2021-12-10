package models

import (
	"github.com/chuakid/cvwo-backend/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int
	Username string     `json:"username" gorm:"unique"`
	Password string     `json:"password"`
	Projects []*Project `gorm:"many2many:UserProjects;"`
}

func (user *User) CreateUser() (int, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, err
	}
	user.Password = string(bytes)

	result := db.DB.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}
