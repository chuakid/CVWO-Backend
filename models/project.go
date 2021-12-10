package models

import (
	"log"
	"strconv"

	"github.com/chuakid/cvwo-backend/db"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID    int     `json:"ID"`
	Name  string  `json:"name"`
	Users []*User `gorm:"many2many:UserProjects"`
	Tasks []*Task
}

func (project *Project) GetProject() error {
	result := db.DB.Preload("Tasks").First(&project)

	return result.Error
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

func (project *Project) GetUsers() ([]User, error) {
	var users []User
	err := db.DB.Model(&project).Association("Users").Find(&users)
	return users, err
}
