package models

import (
	"errors"

	"github.com/chuakid/cvwo-backend/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int        `json:"ID"`
	Username string     `json:"username" gorm:"unique"`
	Password string     `json:"password"`
	Projects []*Project `gorm:"many2many:UserProjects;"`
}

type APIUser struct {
	ID       int
	Username string     `json:"username"`
	Projects []*Project `gorm:"many2many:UserProjects;"`
}

func (user *User) UserExists() bool {
	err := db.DB.Where("username = ?", user.Username).First(&user).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (user *User) GetProjects() ([]APIProjectSummary, error) {
	var projects []APIProjectSummary
	err := db.DB.Model(&user).Association("Projects").Find(&projects)
	return projects, err
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
