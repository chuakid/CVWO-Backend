package models

import (
	"log"
	"strconv"

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

type Project struct {
	gorm.Model
	ID    int
	Name  string  `json:"name"`
	Users []*User `gorm:"many2many:UserProjects"`
}

func (user *User) GetProjects() ([]Project, error) {
	var projects []Project
	err := db.DB.Model(&user).Association("Projects").Find(&projects)
	return projects, err
}

func (project *Project) CreateProject(userid string) (int, error) {
	id, err := strconv.Atoi(userid)
	if err != nil {
		log.Print("Error creating project:", err)
		return 0, err
	}
	user := User{
		ID: id,
	}
	project.Users = append(project.Users, &user)
	result := db.DB.Create(&project)
	if result.Error != nil {
		return 0, result.Error
	}

	return project.ID, nil
}
