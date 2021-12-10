package models

import (
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
	result := db.DB.Preload("Tasks").Preload("Users").First(&project)
	return result.Error
}

func (project *Project) CreateProject(userid string) error {
	id, err := strconv.Atoi(userid)
	if err != nil {
		return err
	}
	user := User{
		ID: id,
	}
	project.Users = append(project.Users, &user)
	result := db.DB.Create(&project)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (project *Project) DeleteProject() error {
	return db.DB.Delete(&project).Error
}

func (project *Project) GetUsers() ([]User, error) {
	var users []User
	err := db.DB.Model(&project).Association("Users").Find(&users)
	return users, err
}
